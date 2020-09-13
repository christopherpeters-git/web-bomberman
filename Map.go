package main

import (
	"container/list"
	"time"
)

type ItemType int

// -1 doesnt work
var globalBombCount uint64 = 0

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
	if f.Contains[0] != nil {
		f.Contains[1] = b
	} else {
		f.Contains[0] = b
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
	ID     uint64
	Owner  *Bomberman
	Time   int
	Radius int
}

func NewBomb(b *Bomberman) Bomb {
	globalBombCount++
	return Bomb{
		ID:     globalBombCount,
		Owner:  b,
		Time:   b.bombTime,
		Radius: b.BombRadius,
	}
}

type Item struct {
	Type ItemType
}

type Wall struct {
	Destructible bool
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

func (b *Bomb) startBomb(x int, y int) {
	time.Sleep(time.Duration(b.Time) * time.Second)
	GameMap.Fields[x][y].explosion()
	for i := 1; i < b.Radius; i++ {
		xPos := x + i
		xNeg := x - i
		yPos := y + i
		yNeg := y - i
		if xPos <= len(GameMap.Fields) {
			GameMap.Fields[xPos][y].explosion()
		}
		if xNeg >= 0 {
			GameMap.Fields[xNeg][y].explosion()
		}
		if yPos <= len(GameMap.Fields[x]) {
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
		element.Value.(*Bomberman).isAlive = false
		for element.Next() != nil {
			element = element.Next()
			element.Value.(*Bomberman).isAlive = false
		}
	}
	for i := 0; i < 2; i++ {
		if f.Contains[i].isDestructible() {
			f.Contains[i] = nil
		}
	}
}
