package main

import (
	"errors"
	"image"
	"image/png"
	"io"
	"os"
	"strconv"
)

/*
Maps a FieldObject to a RGBA color.
*/
var PIXEL_WALL_SOLID = newPixel(0, 0, 0, 255)
var PIXEL_WALL_WEAK = newPixel(66, 65, 66, 255)
var PIXEL_ITEM_BOOST = newPixel(0, 230, 255, 255)
var PIXEL_ITEM_SLOW = newPixel(255, 115, 0, 255)
var PIXEL_ITEM_GHOST = newPixel(0, 26, 255, 255)

/*
Gets a Map and a Path to an Image, which has the same Pixel-Dimensions as the Field-Array.

Loops over all Pixels and adds a FieldObject to the Map according to the RGBA-Value of the Color.
*/
func CreateMapFromImage(m Map, imagePath string) error {
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	file, err := os.Open(imagePath)
	if err != nil {
		return err
	}
	defer file.Close()

	image, err := png.Decode(file)
	if err != nil {
		return err
	}
	if image.Bounds().Dx() != MAP_SIZE || image.Bounds().Dy() != MAP_SIZE {
		return errors.New("Creating map from image failed: png needs to have the height " + strconv.Itoa(MAP_SIZE) + " and width " + strconv.Itoa(MAP_SIZE))
	}

	file, err = os.Open(imagePath)
	if err != nil {
		return err
	}

	pixels, err := getPixels(file)
	if err != nil {
		return err
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

	//j und i vertauscht?
	for i := 0; i < len(pixels); i++ {
		for j := 0; j < len(pixels[i]); j++ {
			if pixels[i][j] == PIXEL_WALL_SOLID {
				m.Fields[j][i].addWall(wSolid)
			}
			if pixels[i][j] == PIXEL_WALL_WEAK {
				m.Fields[j][i].addWall(wWeak)
			}
			if pixels[i][j] == PIXEL_ITEM_BOOST {
				m.Fields[j][i].addItem(&i0)
			}
			if pixels[i][j] == PIXEL_ITEM_SLOW {
				m.Fields[j][i].addItem(&i1)
			}
			if pixels[i][j] == PIXEL_ITEM_GHOST {
				m.Fields[j][i].addItem(&i2)
			}
		}
	}
	return nil
}

/*
Converts a PNG to a two dimensional Pixel-Array.
*/
func getPixels(file io.Reader) ([][]Pixel, error) {
	img, _, err := image.Decode(file)

	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

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

/*
Converts "img.At(x, y).RGBA()" to a Pixel.
"img.At(x, y).RGBA()" returns four uint32 values, we need a Pixel.
*/
func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}

/*
Represents a Pixel with RGBA values.
*/
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
