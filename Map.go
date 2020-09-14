package main

import (
	"container/list"
	"log"
	"time"
)

type ItemType int

// -1 doesnt work
var globalBombCount uint64 = 0
var globalTestMap Map = NewMap(10)

const (
	ItemTypeUpgrade    ItemType = 0
	ItemTypeDowngrade  ItemType = 1
	ItemTypeShortBoost ItemType = 2
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
	return m
}

type Field struct {
	Contains []FieldType
	Player   *list.List
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

func NewField() Field {
	return Field{
		Contains: make([]FieldType, 2),
		Player:   list.New(),
	}
}

type FieldType interface {
	isAccessible() bool
	startEvent()
	isDestructible() bool
}

type Bomb struct {
	ID        uint64
	Owner     *Bomberman
	PositionX int
	PositionY int
	Time      int
	Radius    int
}

//todo *Bomb needed?
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

type Item struct {
	Type ItemType
}

type Wall struct {
	Destructible bool
}

func newWall(destructible bool) *Wall {
	return &Wall{Destructible: destructible}
}

func (b *Bomb) isAccessible() bool {
	return false
}

func (i *Item) isAccessible() bool {
	return true
}
func (w *Wall) isAccessible() bool {
	return false
}

func (b *Bomb) startEvent() {

}
func (i *Item) startEvent() {

}
func (w *Wall) startEvent() {

}

func (b *Bomb) isDestructible() bool {
	return false
}

func (i *Item) isDestructible() bool {
	return false
}

func (w *Wall) isDestructible() bool {
	return w.Destructible
}

func (b *Bomb) startBomb() {
	log.Println("Starting bomb...")
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

func FillTestMap(m Map) {
	w0 := newWall(true)
	w1 := newWall(true)
	w2 := newWall(true)
	w3 := newWall(true)
	w4 := newWall(true)
	w5 := newWall(true)
	w6 := newWall(true)
	w7 := newWall(true)
	w8 := newWall(true)
	w9 := newWall(true)
	w10 := newWall(false)
	w11 := newWall(false)
	w12 := newWall(false)
	w13 := newWall(false)
	w14 := newWall(false)
	w15 := newWall(false)
	w16 := newWall(false)
	w17 := newWall(false)
	w18 := newWall(false)
	w19 := newWall(false)
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
