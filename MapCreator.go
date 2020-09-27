package main

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
)

func CreateMapFromImage(m Map, imagePath string) {

	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)

	file, err := os.Open(imagePath)

	if err != nil {
		fmt.Println("Error: File could not be opened")
		os.Exit(1)
	}

	defer file.Close()

	pixels, err := getPixels(file)

	if err != nil {
		fmt.Println("Error: Image could not be decoded")
		os.Exit(1)
	}
	wSolid := NewWall(false)
	wWeak := NewWall(true)
	i0 := NewItem(FieldObjectItemBoost)
	i1 := NewItem(FieldObjectItemSlow)
	i2 := NewItem(FieldObjectItemGhost)
	p0 := NewPortal(newPosition(9, 3), newPosition(8, 8))
	p1 := NewPortal(newPosition(10, 3), newPosition(11, 8))
	p2 := NewPortal(newPosition(8, 11), newPosition(9, 16))
	p3 := NewPortal(newPosition(11, 11), newPosition(10, 16))
	m.addPortal(&p0)
	m.addPortal(&p1)
	m.addPortal(&p2)
	m.addPortal(&p3)
	//fmt.Println(pixels)

	wallPixel := newPixel(0, 0, 0, 255)

	//j und i vertauscht?
	for i := 0; i < len(pixels); i++ {
		for j := 0; j < len(pixels[i]); j++ {
			if pixels[i][j] == wallPixel {
				m.Fields[j][i].addWall(wSolid)
			}
			if pixels[i][j].R == 66 && pixels[i][j].G == 65 && pixels[i][j].B == 66 && pixels[i][j].A == 255 {
				m.Fields[j][i].addWall(wWeak)
			}
			if pixels[i][j].R == 255 && pixels[i][j].G == 115 && pixels[i][j].B == 0 && pixels[i][j].A == 255 {
				m.Fields[j][i].addItem(&i1)
			}
			if pixels[i][j].R == 0 && pixels[i][j].G == 230 && pixels[i][j].B == 255 && pixels[i][j].A == 255 {
				m.Fields[j][i].addItem(&i0)
			}
			if pixels[i][j].R == 0 && pixels[i][j].G == 26 && pixels[i][j].B == 255 && pixels[i][j].A == 255 {
				m.Fields[j][i].addItem(&i2)
			}
		}
	}
}

// Get the bi-dimensional pixel array
func getPixels(file io.Reader) ([][]Pixel, error) {
	img, _, err := image.Decode(file)

	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	//Überprüfen ob Bild größe der Mapsize entspricht

	var pixels [][]Pixel
	for y := 0; y < height; y++ {
		var row []Pixel
		for x := 0; x < width; x++ {
			row = append(row, rgbaToPixel(img.At(x, y).RGBA()))
		}
		pixels = append(pixels, row)
	}

	return pixels, nil
}

// img.At(x, y).RGBA() returns four uint32 values; we want a Pixel
func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}

// Pixel struct example
type Pixel struct {
	R int
	G int
	B int
	A int
}

func newPixel(r int, g int, b int, a int) Pixel {
	return Pixel{
		R: r,
		G: g,
		B: b,
		A: a,
	}
}
