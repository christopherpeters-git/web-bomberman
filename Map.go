package main

type ItemType int

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
	for i := range m.Fields {
		m.Fields[i] = make([]Field, size)
	}
	return m
}

type Field struct {
	Contains []FieldType
	Player   []*Bomberman
}

func (f *Field) addBomb(b *Bomb) {
	f.Contains = append(f.Contains, b)
}

func NewField() Field {
	return Field{
		Contains: make([]FieldType, 0),
		Player:   make([]*Bomberman, 0),
	}
}

type FieldType interface {
	isAccessible() bool
	startEvent()
}

type Bomb struct {
	Owner  *Bomberman
	Time   int
	Radius int
}

func NewBomb(b *Bomberman) Bomb {
	return Bomb{
		Owner:  b,
		Time:   b.bombTime,
		Radius: b.BombRadius,
	}
}

type Item struct {
	Type ItemType
}

type Wall struct {
	Destructable bool
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

func (b *Bomb) startBomb() {

}
