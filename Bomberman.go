package main

import (
	"log"
	"strconv"
	"time"
)

type Bomberman struct {
	UserID         uint64
	PositionX      int
	PositionY      int
	Name           string
	oldPositionX   int
	oldPositionY   int
	lastBombPlaced time.Time
	BombRadius     int
	bombTime       int
	IsAlive        bool
	IsHit          bool
	ItemActive     bool
	IsMoving       bool
	GhostActive    bool
	hasTeleported  bool
	PlayerReady    bool
	topRightPos    Position
	topLeftPos     Position
	bottomRightPos Position
	bottomLeftPos  Position
	stepMult       float32
	DirUp          bool
	DirDown        bool
	DirLeft        bool
	DirRight       bool
}

func (r *Bomberman) String() string {
	return "Bomberman: {" + strconv.FormatUint(r.UserID, 10) + " | " + strconv.FormatInt(int64(r.PositionX), 10) + " | " + strconv.FormatInt(int64(r.PositionY), 10) + " | " + r.lastBombPlaced.String() + "}"
}

func NewBomberman(userID uint64, positionX int, positionY int, name string) *Bomberman {
	return &Bomberman{
		UserID:         userID,
		PositionX:      positionX,
		PositionY:      positionY,
		oldPositionX:   positionX,
		oldPositionY:   positionY,
		Name:           name,
		BombRadius:     2,
		bombTime:       3,
		IsAlive:        true,
		IsHit:          false,
		ItemActive:     false,
		IsMoving:       false,
		GhostActive:    false,
		hasTeleported:  false,
		PlayerReady:    false,
		topRightPos:    newPosition(positionX+43, positionY+7),
		topLeftPos:     newPosition(positionX+7, positionY+7),
		bottomRightPos: newPosition(positionX+7, positionY+43),
		bottomLeftPos:  newPosition(positionX+43, positionY+43),
		stepMult:       1,
		DirUp:          false,
		DirDown:        false,
		DirLeft:        false,
		DirRight:       false,
	}
}

func (r *Bomberman) placeBomb() {
	bomb := NewBomb(r)
	GameMap.Fields[bomb.PositionX][bomb.PositionY].addBomb(&bomb)
	bomb.startBomb()
}

func (r *Bomberman) teleportTo(xDestination int, yDestination int, currentX int, currentY int) {
	if xDestination < 0 || yDestination < 0 || xDestination > (len(GameMap.Fields)-1) || yDestination > (len(GameMap.Fields[xDestination/FIELD_SIZE])-1) {
		return
		log.Println("Teleport cancelled. Out of Bounds!")
	}
	accessOne := true
	accessTwo := true

	if GameMap.Fields[xDestination][yDestination].Contains[0] != nil {
		accessOne = GameMap.Fields[xDestination][yDestination].Contains[0].isAccessible()
	} else if GameMap.Fields[xDestination][yDestination].Contains[1] != nil {
		accessTwo = GameMap.Fields[xDestination][yDestination].Contains[1].isAccessible()
	}

	if accessOne && accessTwo {
		removePlayerFromList(GameMap.Fields[currentX][currentY].Player, r)
		GameMap.Fields[xDestination][yDestination].Player.PushBack(r)
		r.PositionX = xDestination * FIELD_SIZE
		r.PositionY = yDestination * FIELD_SIZE
		r.oldPositionX = xDestination * FIELD_SIZE
		r.oldPositionY = yDestination * FIELD_SIZE
		r.topRightPos.x = r.PositionX + 43
		r.topRightPos.y = r.PositionY + 7
		r.topLeftPos.x = r.PositionX + 7
		r.topLeftPos.y = r.PositionY + 7
		r.bottomRightPos.x = r.PositionX + 7
		r.bottomRightPos.y = r.PositionY + 43
		r.bottomLeftPos.x = r.PositionX + 43
		r.bottomLeftPos.y = r.PositionY + 43
		//MapChanged()
		r.hasTeleported = true
		//checking for poison
		arrX := pixToArr(r.PositionX)
		arrY := pixToArr(r.PositionY)
		if GameMap.Fields[arrX][arrY].Contains != nil {
			if GameMap.Fields[arrX][arrY].Contains[0] != nil && GameMap.Fields[arrX][arrY].Contains[0].getType() == FieldObjectPoison ||
				GameMap.Fields[arrX][arrY].Contains[1] != nil && GameMap.Fields[arrX][arrY].Contains[1].getType() == FieldObjectPoison {
				r.Kill()
				findWinner()
			}
		}
	}
}

func (r *Bomberman) checkFieldForItem(x int, y int) {
	itemKey := FieldObjectNull
	i := 0

	if GameMap.Fields[x][y].Contains[0] != nil {
		itemKey = GameMap.Fields[x][y].Contains[0].getType()
	} else if GameMap.Fields[x][y].Contains[1] != nil {
		i = 1
		itemKey = GameMap.Fields[x][y].Contains[1].getType()
	}

	if itemKey == 12 {
		portal := GameMap.Fields[x][y].Contains[1].(*Portal)

		if x == portal.portalOne.x && y == portal.portalOne.y {
			r.teleportTo(portal.portalTwo.x, portal.portalTwo.y, x, y)
		} else if x == portal.portalTwo.x && y == portal.portalTwo.y {
			r.teleportTo(portal.portalOne.x, portal.portalOne.y, x, y)
		}
		return
	}

	if r.ItemActive {
		return
	}

	switch itemKey {
	case 0:
		return
	case 6:
		r.ItemActive = true
		r.stepMult = 1.5
		GameMap.Fields[x][y].Contains[i] = nil
		time.AfterFunc(7*time.Second, func() {
			if r.IsAlive {
				r.stepMult = 1
			}
			r.ItemActive = false
		})
	case 7:
		r.ItemActive = true
		r.stepMult = 0.8
		GameMap.Fields[x][y].Contains[i] = nil
		time.AfterFunc(7*time.Second, func() {
			if r.IsAlive {
				r.stepMult = 1
			}
			r.ItemActive = false
		})
	case 8:
		r.ItemActive = true
		r.GhostActive = true
		GameMap.Fields[x][y].Contains[i] = nil

		time.AfterFunc(5*time.Second, func() {
			if r.IsAlive {
				r.GhostActive = false
			}
			r.ItemActive = false
		})
	case 13:
		r.Kill()
		findWinner()
	default:
		return

	}
	//MapChanged()
}

func (b *Bomberman) Reset(x int, y int) {
	b.IsAlive = true
	b.GhostActive = false
	b.ItemActive = false
	b.DirUp = false
	b.DirDown = false
	b.DirLeft = false
	b.DirRight = false
	b.IsHit = false
	b.PlayerReady = false
	b.IsMoving = false
	b.stepMult = STANDARD_STEP_MULTIPLICATOR
	b.BombRadius = STANDARD_BOMB_RADIUS
	b.bombTime = STANDARD_BOMB_TIME
	b.teleportTo(x, y, pixToArr(b.PositionX), pixToArr(b.PositionY))
}

func (b *Bomberman) Kill() {
	b.IsAlive = false
	b.GhostActive = true
	b.ItemActive = true
	b.stepMult = DEATH_STEP_MULTIPLICATOR
}
