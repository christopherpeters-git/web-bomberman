package main

import (
	"container/list"
	"log"
	"time"
)

type ItemType int
type FieldObject int

// -1 doesnt work
var globalBombCount uint64 = 0
var globalTestMap Map = NewMap(10)

const (
	ItemTypeUpgrade    ItemType = 0
	ItemTypeDowngrade  ItemType = 1
	ItemTypeShortBoost ItemType = 2
)

const (
	FieldObjectBomb          FieldObject = 1
	FieldObjectWall          FieldObject = 2
	FieldObjectItemUpgrade   FieldObject = 3
	FieldObjectItemDowngrade FieldObject = 4
	FieldObjectItemBoost     FieldObject = 5
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
	FillTestMap(m)
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
}

func (f *Field) addWall(w *Wall) {
	if f.Contains[0] != nil {
		f.Contains[1] = w
	} else {
		f.Contains[0] = w
	}
}

func (f *Field) explosion() {
	element := f.Player.Front()
	if element != nil {
		element.Value.(*Bomberman).IsAlive = false
		for element.Next() != nil {
			element = element.Next()
			element.Value.(*Bomberman).IsAlive = false
		}
	}
	for i := 0; i < 2; i++ {
		if f.Contains[i] != nil {
			if f.Contains[i].isDestructible() {
				f.Contains[i] = nil
			}
		}
	}
}

type FieldType interface {
	isAccessible() bool
	startEvent()
	isDestructible() bool
	getType() FieldObject
}

type Bomb struct {
	ID        uint64
	Owner     *Bomberman
	PositionX int
	PositionY int
	Time      int
	Radius    int
}

func NewBomb(b *Bomberman) Bomb {
	globalBombCount++
	return Bomb{
		ID:        globalBombCount,
		Owner:     b,
		PositionX: b.PositionX / FIELD_SIZE,
		PositionY: b.PositionY / FIELD_SIZE,
		Time:      b.bombTime,
		Radius:    b.BombRadius,
	}
}

func (b *Bomb) isAccessible() bool {
	return false
}
func (b *Bomb) startEvent() {

}
func (b *Bomb) isDestructible() bool {
	return false
}
func (b *Bomb) getType() FieldObject {
	return FieldObjectBomb
}

func (b *Bomb) startBomb() {
	time.Sleep(time.Duration(b.Time) * time.Second)
	x := b.PositionX
	y := b.PositionY
	GameMap.Fields[x][y].explosion()
	for i := 1; i < b.Radius; i++ {
		xPos := x + i
		xNeg := x - i
		yPos := y + i
		yNeg := y - i
		if xPos < len(GameMap.Fields) {
			GameMap.Fields[xPos][y].explosion()
		}
		if xNeg >= 0 {
			GameMap.Fields[xNeg][y].explosion()
		}
		if yPos < len(GameMap.Fields[x]) {
			GameMap.Fields[x][yPos].explosion()
		}
		if yNeg >= 0 {
			GameMap.Fields[x][yNeg].explosion()
		}
	}
	if GameMap.Fields[x][y].Contains[0] == b {
		GameMap.Fields[x][y].Contains[0] = nil
	} else if GameMap.Fields[x][y].Contains[1] == b {
		GameMap.Fields[x][y].Contains[1] = nil
	}
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
func (i *Item) startEvent() {

}
func (i *Item) isDestructible() bool {
	return false
}
func (i *Item) getType() FieldObject {
	return i.Type
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
func (w *Wall) startEvent() {

}
func (w *Wall) isDestructible() bool {
	return w.Destructible
}
func (w *Wall) getType() FieldObject {
	return FieldObjectWall
}

func FillTestMap(m Map) {
	w0 := NewWall(true)
	w1 := NewWall(true)
	w2 := NewWall(true)
	w3 := NewWall(true)
	w4 := NewWall(true)
	w5 := NewWall(true)
	w6 := NewWall(true)
	w7 := NewWall(true)
	w8 := NewWall(true)
	w9 := NewWall(true)
	w10 := NewWall(false)
	w11 := NewWall(false)
	w12 := NewWall(false)
	w13 := NewWall(false)
	w14 := NewWall(false)
	w15 := NewWall(false)
	w16 := NewWall(false)
	w17 := NewWall(false)
	w18 := NewWall(false)
	w19 := NewWall(false)
	m.Fields[3][0].addWall(w0)
	m.Fields[5][0].addWall(w1)
	m.Fields[2][1].addWall(w10)
	m.Fields[3][1].addWall(w11)
	m.Fields[4][1].addWall(w12)
	m.Fields[6][1].addWall(w13)
	m.Fields[1][2].addWall(w14)
	m.Fields[2][2].addWall(w2)
	m.Fields[4][2].addWall(w15)
	m.Fields[5][2].addWall(w3)
	m.Fields[6][2].addWall(w16)
	m.Fields[1][3].addWall(w4)
	m.Fields[2][3].addWall(w5)
	m.Fields[3][3].addWall(w6)
	m.Fields[0][4].addWall(w7)
	m.Fields[1][4].addWall(w17)
	m.Fields[2][4].addWall(w8)
	m.Fields[4][4].addWall(w9)
	m.Fields[1][5].addWall(w18)
	m.Fields[2][5].addWall(w19)
}
