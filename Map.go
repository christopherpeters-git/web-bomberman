package main

import (
	"container/list"
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"os"
)

type ItemType int
type FieldObject int

var globalBombCount uint64 = 0
var playerDied bool = false

const bombStates = 3

const (
	ItemTypeUpgrade    ItemType = 0
	ItemTypeDowngrade  ItemType = 1
	ItemTypeShortBoost ItemType = 2
)

const (
	FieldObjectNull          FieldObject = 0
	FieldObjectBomb          FieldObject = 1
	FieldObjectBomb1         FieldObject = 10
	FieldObjectBomb2         FieldObject = 11
	FieldObjectWeakWall      FieldObject = 2
	FieldObjectSolidWall     FieldObject = 3
	FieldObjectItemUpgrade   FieldObject = 4
	FieldObjectItemDowngrade FieldObject = 5
	FieldObjectItemBoost     FieldObject = 6
	FieldObjectItemSlow      FieldObject = 7
	FieldObjectItemGhost     FieldObject = 8
	FieldObjectExplosion     FieldObject = 9
	FieldObjectPortal        FieldObject = 12
)

type Map struct {
	Fields [][]Field
}

func NewMap(size int) Map {
	m := Map{Fields: make([][]Field, size)}
	for i := 0; i < len(m.Fields); i++ {
		m.Fields[i] = make([]Field, size)
		for j := 0; j < len(m.Fields[i]); j++ {
			m.Fields[i][j] = NewField()
		}
	}
	CreateMapFromImage(m, "images/map.png")
	return m
}

type Field struct {
	Contains []FieldType
	Player   *list.List
}

func NewField() Field {
	return Field{
		Contains: make([]FieldType, 2),
		Player:   list.New(),
	}
}

func (f *Field) addBomb(b *Bomb) {
	log.Println("added bomb.")
	if f.Contains[0] != nil {
		f.Contains[1] = b
	} else {
		f.Contains[0] = b
	}
	BuildAbstractGameMap()
}

func (f *Field) addWall(w *Wall) {
	if f.Contains[0] != nil {
		f.Contains[1] = w
	} else {
		f.Contains[0] = w
	}
}
func (f *Field) addItem(i *Item) {
	if f.Contains[0] != nil {
		f.Contains[1] = i
	} else {
		f.Contains[0] = i
	}
}

func (m *Map) addPortal(p *Portal) {
	m.Fields[p.portalOne.x][p.portalOne.y].Contains[1] = p
	m.Fields[p.portalTwo.x][p.portalTwo.y].Contains[1] = p
}

func (f *Field) addExplosion(e *Explosion) {
	if f.Contains[0] != nil {
		f.Contains[1] = e
	} else {
		f.Contains[0] = e
	}
}

type eventFunction func(i interface{})

type FieldType interface {
	isAccessible() bool
	startEvent(f eventFunction)
	isDestructible() bool
	getType() FieldObject
}

type Explosion struct {
	ExpFields []Position
}

func newExplosion() Explosion {
	return Explosion{
		ExpFields: make([]Position, 0),
	}
}

func (e *Explosion) isAccessible() bool {
	return true
}

func (e *Explosion) startEvent(f eventFunction) {

}

func (e *Explosion) isDestructible() bool {
	return false
}

func (e *Explosion) getType() FieldObject {
	return FieldObjectExplosion
}

type Item struct {
	Type FieldObject
}

func NewItem(t FieldObject) Item {
	return Item{Type: t}
}

func (i *Item) isAccessible() bool {
	return true
}
func (i *Item) startEvent(f eventFunction) {
}

func (i *Item) isDestructible() bool {
	return false
}
func (i *Item) getType() FieldObject {
	return i.Type
}

type Portal struct {
	iFeelUsed bool
	portalOne Position
	portalTwo Position
}

func NewPortal(portalOne Position, portalTwo Position) Portal {
	return Portal{
		iFeelUsed: false,
		portalOne: portalOne,
		portalTwo: portalTwo,
	}
}

func (p *Portal) isAccessible() bool {
	return true
}

func (p *Portal) startEvent(f eventFunction) {

}

func (p *Portal) isDestructible() bool {
	return false
}

func (p *Portal) getType() FieldObject {
	return FieldObjectPortal
}

type Wall struct {
	Destructible bool
}

func NewWall(destructible bool) *Wall {
	return &Wall{Destructible: destructible}
}

func (w *Wall) isAccessible() bool {
	return false
}
func (w *Wall) startEvent(f eventFunction) {

}
func (w *Wall) isDestructible() bool {
	return w.Destructible
}
func (w *Wall) getType() FieldObject {
	if w.isDestructible() {
		return FieldObjectWeakWall
	} else {
		return FieldObjectSolidWall
	}
}

func BuildAbstractGameMap() {
	//Create map to send
	abstractGameMap = make([][][]FieldObject, len(GameMap.Fields))
	//testToSend := make([][]int, len(GameMap.Fields))
	for i, _ := range GameMap.Fields {
		abstractGameMap[i] = make([][]FieldObject, len(GameMap.Fields[i]))
		//testToSend[i] = make([]int, len(GameMap.Fields[i]))
		for j, _ := range GameMap.Fields[i] {
			abstractGameMap[i][j] = make([]FieldObject, len(GameMap.Fields[i][j].Contains))
			if GameMap.Fields[i][j].Player.Front() != nil {
				//testToSend[i][j] = 1
			}
			for k, _ := range GameMap.Fields[i][j].Contains {
				if GameMap.Fields[i][j].Contains[k] != nil {
					abstractGameMap[i][j][k] = GameMap.Fields[i][j].Contains[k].getType()
				}
			}
		}
	}
}

func CreateMapFromImage(m Map, imagePfad string) {

	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)

	file, err := os.Open(imagePfad)

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

func (m *Map) clearMap() {
	for i := 0; i < len(m.Fields); i++ {
		for j := 0; j < len(m.Fields[0]); j++ {
			m.Fields[i][j] = NewField()
		}
	}
}
