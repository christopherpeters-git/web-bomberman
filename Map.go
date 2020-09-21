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

//var globalTestMap Map = NewMap(10)

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

func (p *Portal) startEvent() {

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
	p0 := NewPortal(newPosition(0, 1), newPosition(0, 18))
	m.addPortal(&p0)

	m.Fields[9][0].addWall(wSolid)
	m.Fields[10][0].addWall(wSolid)

	m.Fields[4][1].addWall(wWeak)
	m.Fields[5][1].addWall(wSolid)
	m.Fields[6][1].addWall(wSolid)
	m.Fields[7][1].addWall(wSolid)
	m.Fields[8][1].addWall(wSolid)
	m.Fields[9][1].addWall(wSolid)
	m.Fields[10][1].addWall(wSolid)
	m.Fields[11][1].addWall(wSolid)
	m.Fields[12][1].addWall(wSolid)
	m.Fields[13][1].addWall(wSolid)
	m.Fields[14][1].addWall(wSolid)
	m.Fields[15][1].addWall(wWeak)

	m.Fields[3][2].addWall(wWeak)
	m.Fields[4][2].addWall(wWeak)
	m.Fields[5][2].addWall(wWeak)
	m.Fields[6][2].addWall(wWeak)
	m.Fields[7][2].addWall(wWeak)
	m.Fields[8][2].addWall(wWeak)
	m.Fields[9][2].addWall(wWeak)
	m.Fields[10][2].addWall(wWeak)
	m.Fields[11][2].addWall(wWeak)
	m.Fields[12][2].addWall(wWeak)
	m.Fields[13][2].addWall(wWeak)
	m.Fields[14][2].addWall(wWeak)
	m.Fields[15][2].addWall(wWeak)
	m.Fields[16][2].addWall(wWeak)

	m.Fields[2][3].addWall(wSolid)
	m.Fields[3][3].addWall(wSolid)
	m.Fields[6][3].addWall(wWeak)
	m.Fields[7][3].addWall(wWeak)
	m.Fields[8][3].addWall(wWeak)
	m.Fields[11][3].addWall(wWeak)
	m.Fields[12][3].addWall(wWeak)
	m.Fields[13][3].addWall(wWeak)
	m.Fields[16][3].addWall(wSolid)
	m.Fields[17][3].addWall(wSolid)

	m.Fields[2][4].addWall(wSolid)
	m.Fields[3][4].addWall(wWeak)
	m.Fields[4][4].addWall(wWeak)
	m.Fields[5][4].addWall(wWeak)
	m.Fields[6][4].addWall(wWeak)
	m.Fields[7][4].addWall(wWeak)
	m.Fields[8][4].addWall(wWeak)
	m.Fields[9][4].addWall(wWeak)
	m.Fields[10][4].addWall(wWeak)
	m.Fields[11][4].addWall(wWeak)
	m.Fields[12][4].addWall(wWeak)
	m.Fields[13][4].addWall(wWeak)
	m.Fields[14][4].addWall(wWeak)
	m.Fields[15][4].addWall(wWeak)
	m.Fields[16][4].addWall(wWeak)
	m.Fields[17][4].addWall(wSolid)

	m.Fields[1][5].addWall(wWeak)
	m.Fields[2][5].addWall(wWeak)
	m.Fields[3][5].addWall(wWeak)
	m.Fields[4][5].addWall(wSolid)
	m.Fields[5][5].addItem(&i1)
	m.Fields[6][5].addWall(wSolid)
	m.Fields[7][5].addWall(wSolid)
	m.Fields[8][5].addWall(wSolid)
	m.Fields[9][5].addWall(wSolid)
	m.Fields[10][5].addWall(wSolid)
	m.Fields[11][5].addWall(wSolid)
	m.Fields[12][5].addWall(wSolid)
	m.Fields[13][5].addWall(wSolid)
	m.Fields[14][5].addItem(&i1)
	m.Fields[15][5].addWall(wSolid)
	m.Fields[16][5].addWall(wWeak)
	m.Fields[17][5].addWall(wWeak)
	m.Fields[18][5].addWall(wWeak)

	m.Fields[0][6].addWall(wWeak)
	m.Fields[1][6].addWall(wSolid)
	m.Fields[3][6].addWall(wWeak)
	m.Fields[4][6].addItem(&i1)
	m.Fields[6][6].addWall(wWeak)
	m.Fields[9][6].addWall(wWeak)
	m.Fields[10][6].addWall(wWeak)
	m.Fields[13][6].addWall(wWeak)
	m.Fields[15][6].addItem(&i1)
	m.Fields[16][6].addWall(wWeak)
	m.Fields[18][6].addWall(wSolid)
	m.Fields[19][6].addWall(wWeak)

	m.Fields[0][7].addWall(wWeak)
	m.Fields[1][7].addWall(wSolid)
	m.Fields[2][7].addWall(wWeak)
	m.Fields[3][7].addWall(wWeak)
	m.Fields[4][7].addWall(wSolid)
	m.Fields[5][7].addWall(wWeak)
	m.Fields[6][7].addWall(wSolid)
	m.Fields[7][7].addWall(wWeak)
	m.Fields[8][7].addWall(wWeak)
	m.Fields[9][7].addWall(wSolid)
	m.Fields[10][7].addWall(wSolid)
	m.Fields[11][7].addWall(wWeak)
	m.Fields[12][7].addWall(wWeak)
	m.Fields[13][7].addWall(wSolid)
	m.Fields[14][7].addWall(wWeak)
	m.Fields[15][7].addWall(wSolid)
	m.Fields[16][7].addWall(wWeak)
	m.Fields[17][7].addWall(wWeak)
	m.Fields[18][7].addWall(wSolid)
	m.Fields[19][7].addWall(wWeak)

	m.Fields[0][8].addItem(&i0)
	m.Fields[1][8].addWall(wSolid)
	m.Fields[2][8].addWall(wWeak)
	m.Fields[3][8].addWall(wWeak)
	m.Fields[4][8].addWall(wSolid)
	m.Fields[6][8].addWall(wWeak)
	m.Fields[9][8].addItem(&i1)
	m.Fields[10][8].addItem(&i1)
	m.Fields[13][8].addWall(wWeak)
	m.Fields[15][8].addWall(wSolid)
	m.Fields[16][8].addWall(wWeak)
	m.Fields[17][8].addWall(wWeak)
	m.Fields[18][8].addWall(wSolid)
	m.Fields[19][8].addItem(&i0)

	m.Fields[0][9].addWall(wSolid)
	m.Fields[1][9].addWall(wSolid)
	m.Fields[2][9].addItem(&i1)
	m.Fields[3][9].addWall(wWeak)
	m.Fields[4][9].addWall(wSolid)
	m.Fields[5][9].addWall(wWeak)
	m.Fields[6][9].addWall(wWeak)
	m.Fields[7][9].addWall(wSolid)
	m.Fields[8][9].addItem(&i1)
	m.Fields[9][9].addItem(&i2)
	m.Fields[10][9].addItem(&i2)
	m.Fields[11][9].addItem(&i1)
	m.Fields[12][9].addWall(wSolid)
	m.Fields[13][9].addWall(wWeak)
	m.Fields[14][9].addWall(wWeak)
	m.Fields[15][9].addWall(wSolid)
	m.Fields[16][9].addWall(wWeak)
	m.Fields[17][9].addItem(&i1)
	m.Fields[18][9].addWall(wSolid)
	m.Fields[19][9].addWall(wSolid)

	m.Fields[0][10].addWall(wSolid)
	m.Fields[1][10].addWall(wSolid)
	m.Fields[2][10].addItem(&i1)
	m.Fields[3][10].addWall(wWeak)
	m.Fields[4][10].addWall(wSolid)
	m.Fields[5][10].addWall(wWeak)
	m.Fields[6][10].addWall(wWeak)
	m.Fields[7][10].addWall(wSolid)
	m.Fields[8][10].addItem(&i1)
	m.Fields[9][10].addItem(&i2)
	m.Fields[10][10].addItem(&i2)
	m.Fields[11][10].addItem(&i1)
	m.Fields[12][10].addWall(wSolid)
	m.Fields[13][10].addWall(wWeak)
	m.Fields[14][10].addWall(wWeak)
	m.Fields[15][10].addWall(wSolid)
	m.Fields[16][10].addWall(wWeak)
	m.Fields[17][10].addItem(&i1)
	m.Fields[18][10].addWall(wSolid)
	m.Fields[19][10].addWall(wSolid)

	m.Fields[0][11].addItem(&i0)
	m.Fields[1][11].addWall(wSolid)
	m.Fields[2][11].addWall(wWeak)
	m.Fields[3][11].addWall(wWeak)
	m.Fields[4][11].addWall(wSolid)
	m.Fields[6][11].addWall(wWeak)
	m.Fields[9][11].addItem(&i1)
	m.Fields[10][11].addItem(&i1)
	m.Fields[13][11].addWall(wWeak)
	m.Fields[15][11].addWall(wSolid)
	m.Fields[16][11].addWall(wWeak)
	m.Fields[17][11].addWall(wWeak)
	m.Fields[18][11].addWall(wSolid)
	m.Fields[19][11].addItem(&i0)

	m.Fields[0][12].addWall(wWeak)
	m.Fields[1][12].addWall(wSolid)
	m.Fields[2][12].addWall(wWeak)
	m.Fields[3][12].addWall(wWeak)
	m.Fields[4][12].addWall(wSolid)
	m.Fields[5][12].addWall(wWeak)
	m.Fields[6][12].addWall(wSolid)
	m.Fields[7][12].addWall(wWeak)
	m.Fields[8][12].addWall(wWeak)
	m.Fields[9][12].addWall(wSolid)
	m.Fields[10][12].addWall(wSolid)
	m.Fields[11][12].addWall(wWeak)
	m.Fields[12][12].addWall(wWeak)
	m.Fields[13][12].addWall(wSolid)
	m.Fields[14][12].addWall(wWeak)
	m.Fields[15][12].addWall(wSolid)
	m.Fields[16][12].addWall(wWeak)
	m.Fields[17][12].addWall(wWeak)
	m.Fields[18][12].addWall(wSolid)
	m.Fields[19][12].addWall(wWeak)

	m.Fields[0][13].addWall(wWeak)
	m.Fields[1][13].addWall(wSolid)
	m.Fields[3][13].addWall(wWeak)
	m.Fields[4][13].addItem(&i1)
	m.Fields[6][13].addWall(wWeak)
	m.Fields[9][13].addWall(wWeak)
	m.Fields[10][13].addWall(wWeak)
	m.Fields[13][13].addWall(wWeak)
	m.Fields[15][13].addItem(&i1)
	m.Fields[16][13].addWall(wWeak)
	m.Fields[18][13].addWall(wSolid)
	m.Fields[19][13].addWall(wWeak)

	m.Fields[1][14].addWall(wWeak)
	m.Fields[2][14].addWall(wWeak)
	m.Fields[3][14].addWall(wWeak)
	m.Fields[4][14].addWall(wSolid)
	m.Fields[5][14].addItem(&i1)
	m.Fields[6][14].addWall(wSolid)
	m.Fields[7][14].addWall(wSolid)
	m.Fields[8][14].addWall(wSolid)
	m.Fields[9][14].addWall(wSolid)
	m.Fields[10][14].addWall(wSolid)
	m.Fields[11][14].addWall(wSolid)
	m.Fields[12][14].addWall(wSolid)
	m.Fields[13][14].addWall(wSolid)
	m.Fields[14][14].addItem(&i1)
	m.Fields[15][14].addWall(wSolid)
	m.Fields[16][14].addWall(wWeak)
	m.Fields[17][14].addWall(wWeak)
	m.Fields[18][14].addWall(wWeak)

	m.Fields[2][15].addWall(wSolid)
	m.Fields[3][15].addWall(wWeak)
	m.Fields[4][15].addWall(wWeak)
	m.Fields[5][15].addWall(wWeak)
	m.Fields[6][15].addWall(wWeak)
	m.Fields[7][15].addWall(wWeak)
	m.Fields[8][15].addWall(wWeak)
	m.Fields[9][15].addWall(wWeak)
	m.Fields[10][15].addWall(wWeak)
	m.Fields[11][15].addWall(wWeak)
	m.Fields[12][15].addWall(wWeak)
	m.Fields[13][15].addWall(wWeak)
	m.Fields[14][15].addWall(wWeak)
	m.Fields[15][15].addWall(wWeak)
	m.Fields[16][15].addWall(wWeak)
	m.Fields[17][15].addWall(wSolid)

	m.Fields[2][16].addWall(wSolid)
	m.Fields[3][16].addWall(wSolid)
	m.Fields[6][16].addWall(wWeak)
	m.Fields[7][16].addWall(wWeak)
	m.Fields[8][16].addWall(wWeak)
	m.Fields[11][16].addWall(wWeak)
	m.Fields[12][16].addWall(wWeak)
	m.Fields[13][16].addWall(wWeak)
	m.Fields[16][16].addWall(wSolid)
	m.Fields[17][16].addWall(wSolid)

	m.Fields[3][17].addWall(wWeak)
	m.Fields[4][17].addWall(wWeak)
	m.Fields[5][17].addWall(wWeak)
	m.Fields[6][17].addWall(wWeak)
	m.Fields[7][17].addWall(wWeak)
	m.Fields[8][17].addWall(wWeak)
	m.Fields[9][17].addWall(wWeak)
	m.Fields[10][17].addWall(wWeak)
	m.Fields[11][17].addWall(wWeak)
	m.Fields[12][17].addWall(wWeak)
	m.Fields[13][17].addWall(wWeak)
	m.Fields[14][17].addWall(wWeak)
	m.Fields[15][17].addWall(wWeak)
	m.Fields[16][17].addWall(wWeak)

	m.Fields[4][18].addWall(wWeak)
	m.Fields[5][18].addWall(wSolid)
	m.Fields[6][18].addWall(wSolid)
	m.Fields[7][18].addWall(wSolid)
	m.Fields[8][18].addWall(wSolid)
	m.Fields[9][18].addWall(wSolid)
	m.Fields[10][18].addWall(wSolid)
	m.Fields[11][18].addWall(wSolid)
	m.Fields[12][18].addWall(wSolid)
	m.Fields[13][18].addWall(wSolid)
	m.Fields[14][18].addWall(wSolid)
	m.Fields[15][18].addWall(wWeak)

	m.Fields[9][19].addWall(wSolid)
	m.Fields[10][19].addWall(wSolid)
	//rand.Seed(time.Now().UTC().UnixNano())
	// (i != 0 || j != 0) && (i != 19 || j != 19) && (i != 0 || j != 19) && (i != 19 || j != 0)
	//for i := 0; i < len(m.Fields); i++ {
	//	for j := 0; j < len(m.Fields[i]); j++ {
	//		if i != 0 && j != 0 && i != 19 && j != 19 {
	//			random := rand.Intn(5)
	//			if random == 1 {
	//				m.Fields[i][j].addWall(wSolid)
	//			} else if random == 2 {
	//				m.Fields[i][j].addWall(wWeak)
	//			}
	//		}
	//	}
	//}

	//for i := 0; i < len(m.Fields); i++ {
	//	for j := 0; j < len(m.Fields[i]); j++ {
	//		if i != 0 && j != 0 && i != 19 && j != 19 {
	//			random := rand.Intn(45) + 1
	//
	//			if random == 15 {
	//				m.Fields[i][j].addItem(&i0)
	//			} else if random == 30 {
	//				m.Fields[i][j].addItem(&i1)
	//			} else if random == 45 {
	//				m.Fields[i][j].addItem(&i2)
	//			}
	//		}
	//	}
	//}

}
