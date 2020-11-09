package main

import (
	"container/list"
	"log"
)

type FieldObject int

var globalBombCount uint64 = 0
var playerDied = false

const bombStates = 3

/*
FieldObject "Enum"
*/
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
	FieldObjectPoison        FieldObject = 13
)

/*
Represents the Gamemap through a two dimensional Array of Fields.
*/
type Map struct {
	Fields [][]Field
}

/*
Initialises the Map.
Fills all Array-Positions with empty Field and then fills it with CreateMapFromImage.
*/
func NewMap(size int) Map {
	m := Map{Fields: make([][]Field, size)}
	for i := 0; i < len(m.Fields); i++ {
		m.Fields[i] = make([]Field, size)
		for j := 0; j < len(m.Fields[i]); j++ {
			m.Fields[i][j] = NewField()
		}
	}
	if err := CreateMapFromImage(m, "images/testMap.png"); err != nil {
		log.Fatal(err)
	}
	return m
}

/*
Fills the Map with empty Field
*/
func (m *Map) clear() {
	for i := 0; i < len(m.Fields); i++ {
		for j := 0; j < len(m.Fields[0]); j++ {
			m.Fields[i][j] = NewField()
		}
	}
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

/*
Adds a Bomb-Object to the Contains-Array of the Field.
*/
func (f *Field) addBomb(b *Bomb) {
	log.Println("added bomb.")
	if f.Contains[0] != nil {
		f.Contains[1] = b
	} else {
		f.Contains[0] = b
	}
	//MapChanged()
}

/*
Adds a Wall-Object to the Contains-Array of the Field.
*/
func (f *Field) addWall(w *Wall) {
	if f.Contains[0] != nil {
		f.Contains[1] = w
	} else {
		f.Contains[0] = w
	}
}

/*
Adds a Item-Object to the Contains-Array of the Field.
*/
func (f *Field) addItem(i *Item) {
	if f.Contains[0] != nil {
		f.Contains[1] = i
	} else {
		f.Contains[0] = i
	}
}

/*
Adds a Portal-Object to the Contains-Array of the Field.
*/
func (m *Map) addPortal(p *Portal) {
	m.Fields[p.portalOne.x][p.portalOne.y].Contains[1] = p
	m.Fields[p.portalTwo.x][p.portalTwo.y].Contains[1] = p
}

/*
Casts an Explosion on a Field. The Explosion kills al Player on the Field and destroys destructible Walls.
*/
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

/*
Adds a Explosion-Object to the Contains-Array of the Field.
*/
func (f *Field) addExplosion(e *Explosion) {
	if f.Contains[0] != nil {
		f.Contains[1] = e
	} else {
		f.Contains[0] = e
	}
}

/*
Adds a Poison-Object to the Contains-Array of the Field.
*/
func (f *Field) addPoison(p *Poison) {
	if f.Contains[0] != nil {
		f.Contains[1] = p
	} else {
		f.Contains[0] = p
	}
}

type eventFunction func(i interface{})

/*
Represents a FieldType.
*/
type FieldType interface {
	isAccessible() bool
	startEvent(f eventFunction)
	isDestructible() bool
	getType() FieldObject
}

type Poison struct {
}

func newPoison() Poison {
	return Poison{}
}

func (p *Poison) isAccessible() bool {
	return true
}

func (p *Poison) startEvent(f eventFunction) {
}

func (p *Poison) isDestructible() bool {
	return false
}

func (p *Poison) getType() FieldObject {
	return FieldObjectPoison
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

/*
Represents the Fields of the Map without Players, only FieldObjects.
*/
func BuildAbstractGameMap() [][][]FieldObject {
	//Create map to send
	newAbstractMap := make([][][]FieldObject, len(GameMap.Fields))
	//testToSend := make([][]int, len(GameMap.Fields))
	for i, _ := range GameMap.Fields {
		newAbstractMap[i] = make([][]FieldObject, len(GameMap.Fields[i]))
		//testToSend[i] = make([]int, len(GameMap.Fields[i]))
		for j, _ := range GameMap.Fields[i] {
			newAbstractMap[i][j] = make([]FieldObject, len(GameMap.Fields[i][j].Contains))
			if GameMap.Fields[i][j].Player.Front() != nil {
				//testToSend[i][j] = 1
			}
			for k, _ := range GameMap.Fields[i][j].Contains {
				if GameMap.Fields[i][j].Contains[k] != nil {
					newAbstractMap[i][j][k] = GameMap.Fields[i][j].Contains[k].getType()
				}
			}
		}
	}
	return newAbstractMap
}
