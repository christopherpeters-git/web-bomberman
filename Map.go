package main

import (
	"container/list"
	"log"
	"math/rand"
	"time"
)

type ItemType int
type FieldObject int

// -1 doesnt work
var globalBombCount uint64 = 0
var globalTestMap Map = NewMap(10)

const bombStates = 3

const (
	ItemTypeUpgrade    ItemType = 0
	ItemTypeDowngrade  ItemType = 1
	ItemTypeShortBoost ItemType = 2
)

const (
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

func (f *Field) addExplosion(e *Explosion) {
	if f.Contains[0] != nil {
		f.Contains[1] = e
	} else {
		f.Contains[0] = e
	}
}

func (f *Field) explosion() bool {
	element := f.Player.Front()
	if element != nil {
		element.Value.(*Bomberman).IsAlive = false
		element.Value.(*Bomberman).GhostActive = true
		element.Value.(*Bomberman).stepMult = 0.5
		for element.Next() != nil {
			element = element.Next()
			element.Value.(*Bomberman).IsAlive = false
			element.Value.(*Bomberman).GhostActive = true
			element.Value.(*Bomberman).stepMult = 0.5
		}
	}
	for i := 0; i < 2; i++ {
		if f.Contains[i] != nil {
			if f.Contains[i].isDestructible() {
				f.Contains[i] = nil

			} else {
				return true
			}
		}
	}
	return false
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
	state     int
}

func NewBomb(b *Bomberman) Bomb {
	globalBombCount++
	return Bomb{
		ID:        globalBombCount,
		Owner:     b,
		PositionX: (b.PositionX + FIELD_SIZE/2) / FIELD_SIZE,
		PositionY: (b.PositionY + FIELD_SIZE/2) / FIELD_SIZE,
		Time:      b.bombTime,
		Radius:    b.BombRadius,
		state:     0,
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
	if b.state == 0 {
		return FieldObjectBomb1
	} else if b.state == 1 {
		return FieldObjectBomb2
	} else {
		return FieldObjectBomb
	}

}

func (b *Bomb) startBomb() {
	//Change to Loop
	time.Sleep((time.Duration(b.Time) / bombStates) * time.Second)
	b.state++
	BuildAbstractGameMap()

	time.Sleep((time.Duration(b.Time) / bombStates) * time.Second)
	b.state++
	BuildAbstractGameMap()

	time.Sleep((time.Duration(b.Time) / bombStates) * time.Second)
	b.state = 0

	e := newExplosion()
	x := b.PositionX
	y := b.PositionY
	xPosHitSolidWall, xNegHitSolidWall, yPosHitSolidWall, yNegHitSolidWall := false, false, false, false
	GameMap.Fields[x][y].explosion()
	e.ExpFields = append(e.ExpFields, newPosition(x, y))
	GameMap.Fields[x][y].addExplosion(&e)
	for i := 1; i < b.Radius; i++ {
		xPos := x + i
		xNeg := x - i
		yPos := y + i
		yNeg := y - i
		if xPos < len(GameMap.Fields) {
			if !xPosHitSolidWall {
				xPosHitSolidWall = GameMap.Fields[xPos][y].explosion()
				if !xPosHitSolidWall {
					e.ExpFields = append(e.ExpFields, newPosition(xPos, y))
					GameMap.Fields[xPos][y].addExplosion(&e)
				}
			}
		}
		if xNeg >= 0 {
			if !xNegHitSolidWall {
				xNegHitSolidWall = GameMap.Fields[xNeg][y].explosion()
				if !xNegHitSolidWall {
					e.ExpFields = append(e.ExpFields, newPosition(xNeg, y))
					GameMap.Fields[xNeg][y].addExplosion(&e)
				}
			}
		}
		if yPos < len(GameMap.Fields[x]) {
			if !yPosHitSolidWall {
				yPosHitSolidWall = GameMap.Fields[x][yPos].explosion()
				if !yPosHitSolidWall {
					e.ExpFields = append(e.ExpFields, newPosition(x, yPos))
					GameMap.Fields[x][yPos].addExplosion(&e)
				}
			}
		}
		if yNeg >= 0 {
			if !yNegHitSolidWall {
				yNegHitSolidWall = GameMap.Fields[x][yNeg].explosion()
				if !yNegHitSolidWall {
					e.ExpFields = append(e.ExpFields, newPosition(x, yNeg))
					GameMap.Fields[x][yNeg].addExplosion(&e)
				}
			}
		}
	}
	if GameMap.Fields[x][y].Contains[0] == b {
		GameMap.Fields[x][y].Contains[0] = nil
	} else if GameMap.Fields[x][y].Contains[1] == b {
		GameMap.Fields[x][y].Contains[1] = nil
	}
	BuildAbstractGameMap()
	time.Sleep(900 * time.Millisecond)
	for i := 0; i < len(e.ExpFields); i++ {
		if GameMap.Fields[e.ExpFields[i].x][e.ExpFields[i].y].Contains[0] != nil {
			if GameMap.Fields[e.ExpFields[i].x][e.ExpFields[i].y].Contains[0].getType() == 9 {
				GameMap.Fields[e.ExpFields[i].x][e.ExpFields[i].y].Contains[0] = nil
			}
		}
		if GameMap.Fields[e.ExpFields[i].x][e.ExpFields[i].y].Contains[1] != nil {
			if GameMap.Fields[e.ExpFields[i].x][e.ExpFields[i].y].Contains[1].getType() == 9 {
				GameMap.Fields[e.ExpFields[i].x][e.ExpFields[i].y].Contains[1] = nil
			}
		}
	}
	BuildAbstractGameMap()
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

func (e *Explosion) startEvent() {

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
	if w.isDestructible() {
		return FieldObjectWeakWall
	} else {
		return FieldObjectSolidWall
	}
}

func FillTestMap(m Map) {
	//w0 := NewWall(true)
	//w1 := NewWall(true)
	//w2 := NewWall(true)
	//w3 := NewWall(true)
	//w4 := NewWall(true)
	//w5 := NewWall(true)
	//w6 := NewWall(true)
	//w7 := NewWall(true)
	//w8 := NewWall(true)
	//w9 := NewWall(true)
	//w10 := NewWall(false)
	//w11 := NewWall(false)
	//w12 := NewWall(false)
	//w13 := NewWall(false)
	//w14 := NewWall(false)
	//w15 := NewWall(false)
	//w16 := NewWall(false)
	//w17 := NewWall(false)
	//w18 := NewWall(false)
	//w19 := NewWall(false)
	//i0 := NewItem(FieldObjectItemBoost)
	//i1 := NewItem(FieldObjectItemSlow)
	//i2 := NewItem(FieldObjectItemGhost)
	//m.Fields[3][0].addWall(w0)
	//m.Fields[5][0].addWall(w1)
	//m.Fields[2][1].addWall(w10)
	//m.Fields[3][1].addWall(w11)
	//m.Fields[4][1].addWall(w12)
	//m.Fields[6][1].addWall(w13)
	//m.Fields[1][2].addWall(w14)
	//m.Fields[2][2].addWall(w2)
	//m.Fields[4][2].addWall(w15)
	//m.Fields[5][2].addWall(w3)
	//m.Fields[6][2].addWall(w16)
	//m.Fields[1][3].addWall(w4)
	//m.Fields[2][3].addWall(w5)
	//m.Fields[3][3].addWall(w6)
	//m.Fields[0][4].addWall(w7)
	//m.Fields[1][4].addWall(w17)
	//m.Fields[2][4].addWall(w8)
	//m.Fields[4][4].addWall(w9)
	//m.Fields[1][5].addWall(w18)
	//m.Fields[2][5].addWall(w19)
	//m.Fields[8][8].addItem(&i0)
	//m.Fields[8][6].addItem(&i1)
	//m.Fields[8][5].addItem(&i2)

	wSolid := NewWall(false)
	wWeak := NewWall(true)
	i0 := NewItem(FieldObjectItemBoost)
	i1 := NewItem(FieldObjectItemSlow)
	i2 := NewItem(FieldObjectItemGhost)
	rand.Seed(time.Now().UTC().UnixNano())
	// (i != 0 || j != 0) && (i != 19 || j != 19) && (i != 0 || j != 19) && (i != 19 || j != 0)
	for i := 0; i < len(m.Fields); i++ {
		for j := 0; j < len(m.Fields[i]); j++ {
			if i != 0 && j != 0 && i != 19 && j != 19 {
				random := rand.Intn(5)
				if random == 1 {
					m.Fields[i][j].addWall(wSolid)
				} else if random == 2 {
					m.Fields[i][j].addWall(wWeak)
				}
			}
		}
	}

	for i := 0; i < len(m.Fields); i++ {
		for j := 0; j < len(m.Fields[i]); j++ {
			if i != 0 && j != 0 && i != 19 && j != 19 {
				random := rand.Intn(45) + 1

				if random == 15 {
					m.Fields[i][j].addItem(&i0)
				} else if random == 30 {
					m.Fields[i][j].addItem(&i1)
				} else if random == 45 {
					m.Fields[i][j].addItem(&i2)
				}
			}
		}
	}

}
