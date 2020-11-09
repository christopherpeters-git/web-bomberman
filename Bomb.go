package main

import (
	"log"
	"time"
)

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
		PositionX: pixToArr(b.PositionX),
		PositionY: pixToArr(b.PositionY),
		Time:      b.bombTime,
		Radius:    b.BombRadius,
		state:     0,
	}
}

func (b *Bomb) isAccessible() bool {
	return false
}
func (b *Bomb) startEvent(f eventFunction) {

}
func (b *Bomb) isDestructible() bool {
	return false
}
func (b *Bomb) getType() FieldObject {
	if b.state == 0 {
		return FieldObjectBombState1
	} else if b.state == 1 {
		return FieldObjectBombState2
	} else {
		return FieldObjectBomb
	}

}

func (b *Bomb) startBomb() {
	//Change to Loop
	time.Sleep((time.Duration(b.Time) / bombStates) * time.Second)
	b.state++
	//MapChanged()
	time.Sleep((time.Duration(b.Time) / bombStates) * time.Second)
	b.state++
	//MapChanged()

	time.Sleep((time.Duration(b.Time) / bombStates) * time.Second)
	b.state = 0

	e := newExplosion()
	x := b.PositionX
	y := b.PositionY
	xPosHitSolidWall, xNegHitSolidWall, yPosHitSolidWall, yNegHitSolidWall := false, false, false, false
	GameMap.Fields[x][y].explosion()
	e.ExpFields = append(e.ExpFields, newPosition(x, y))
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

	if sessionRunning && playerDied {
		log.Println("check for remaining players")
		playerDied = false
		findWinner()
	}

	if GameMap.Fields[x][y].Contains[0] == b {
		GameMap.Fields[x][y].Contains[0] = nil
	} else if GameMap.Fields[x][y].Contains[1] == b {
		GameMap.Fields[x][y].Contains[1] = nil
	}
	GameMap.Fields[x][y].addExplosion(&e)
	//MapChanged()
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
	//MapChanged()
}
