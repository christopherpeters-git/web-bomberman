package main

import (
	"container/list"
	"log"
)

type FieldObject int

var globalBombCount uint64 = 0
var playerDied = false

const bombStates = 3

const (
	FieldObjectNull          FieldObject = 0
	FieldObjectBomb          FieldObject = 1
	FieldObjectWeakWall      FieldObject = 2
	FieldObjectSolidWall     FieldObject = 3
	FieldObjectItemUpgrade   FieldObject = 4
	FieldObjectItemDowngrade FieldObject = 5
	FieldObjectItemBoost     FieldObject = 6
	FieldObjectItemSlow      FieldObject = 7
	FieldObjectItemGhost     FieldObject = 8
	FieldObjectExplosion     FieldObject = 9
	FieldObjectBombState1    FieldObject = 10
	FieldObjectBombState2    FieldObject = 11
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

	if err := CreateMapFromImage(m, "images/map3.png"); err != nil {
		log.Fatal(err)
	}
	return m
}

func (m *Map) clear() {
	for i := 0; i < len(m.Fields); i++ {
		for j := 0; j < len(m.Fields[0]); j++ {
			m.Fields[i][j] = NewField()
		}
	}
}

func (m *Map) addPortal(p *Portal) {
	m.Fields[p.portalOne.x][p.portalOne.y].Contains[1] = p
	m.Fields[p.portalTwo.x][p.portalTwo.y].Contains[1] = p
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

func (f *Field) explosion() bool {
	killAllPlayersOnField(f.Player)
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
