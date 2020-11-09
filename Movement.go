package main

import (
	"container/list"
	"log"
)

/*
Checks if a Field on the Map can be accessed by the Player.
If a Field is accessible or the Player is in Ghost-Mode, the Old-Position gets updated and true gets returned.
*/
func (b *Bomberman) isFieldAccessible(x int, y int) bool {
	isAccessNull := true
	isAccessOne := true
	arrayPosX := (x + FIELD_SIZE/2) / FIELD_SIZE
	arrayPosY := (y + FIELD_SIZE/2) / FIELD_SIZE
	if GameMap.Fields[arrayPosX][arrayPosY].Contains[0] != nil {
		isAccessNull = GameMap.Fields[arrayPosX][arrayPosY].Contains[0].isAccessible()
	}
	if GameMap.Fields[arrayPosX][arrayPosY].Contains[1] != nil {
		isAccessOne = GameMap.Fields[arrayPosX][arrayPosY].Contains[1].isAccessible()
	}

	isAccessible := isAccessNull && isAccessOne
	if isAccessible || b.GhostActive {
		b.oldPositionX = b.PositionX
		b.oldPositionY = b.PositionY
	}
	return isAccessible
}

/*
Removes a Bomberman from the List.
*/
func removePlayerFromList(l *list.List, b *Bomberman) {
	element := l.Front()
	if element != nil {
		//log.Println(b)
		//log.Println(element.Value.(*Bomberman))
		//log.Println(element.Value.(*Bomberman).UserID == b.UserID)
		if element.Value.(*Bomberman).UserID == b.UserID {
			l.Remove(element)
			return
		}
		for element.Next() != nil {
			element = element.Next()
			if element.Value.(*Bomberman).UserID == b.UserID {
				l.Remove(element)
				return
			}
		}
	}
	log.Println("Player not found in list")
}

/*
Checks if the Movement is in Bounds of the Map. If and Array Position needs to be updated, checks if Field is Accessible and
updates the Player-Position if so.
*/
func (r *Bomberman) moveIfLegal(x int, y int) bool {
	if x < 0 || y < 0 || x > (len(GameMap.Fields)-1)*FIELD_SIZE || y > (len(GameMap.Fields[x/FIELD_SIZE])-1)*FIELD_SIZE {
		return false
	}
	oldPosX := (r.PositionX + FIELD_SIZE/2) / FIELD_SIZE
	oldPosY := (r.PositionY + FIELD_SIZE/2) / FIELD_SIZE
	arrayPosX := (x + FIELD_SIZE/2) / FIELD_SIZE
	arrayPosY := (y + FIELD_SIZE/2) / FIELD_SIZE
	inBounds := arrayPosX >= 0 && arrayPosY >= 0 && arrayPosX < len(GameMap.Fields) && arrayPosY < len(GameMap.Fields[arrayPosX])
	if inBounds {
		if oldPosX != arrayPosX {
			if r.isFieldAccessible(x, y) || r.GhostActive {
				/*
					On Teleport the Player Position already gets updated.
				*/
				if r.hasTeleported {
					r.hasTeleported = false
					return false
				}
				removePlayerFromList(GameMap.Fields[oldPosX][arrayPosY].Player, r)
				GameMap.Fields[arrayPosX][arrayPosY].Player.PushBack(r)
				r.checkFieldForItem(arrayPosX, arrayPosY)

				return true
			} else {
				return false
			}
		} else if oldPosY != arrayPosY {
			if r.isFieldAccessible(x, y) || r.GhostActive {
				if r.hasTeleported {
					r.hasTeleported = false
					return false
				}
				removePlayerFromList(GameMap.Fields[arrayPosX][oldPosY].Player, r)
				GameMap.Fields[arrayPosX][arrayPosY].Player.PushBack(r)
				r.checkFieldForItem(arrayPosX, arrayPosY)
				return true
			} else {
				return false
			}
		}
		r.oldPositionX = r.PositionX
		r.oldPositionY = r.PositionY
		return true
	}

	return false
}

/*
???
*/
func outerEdges(x int, y int) bool {
	if x < 0 || y < 0 || x > (len(GameMap.Fields))*FIELD_SIZE || y > (len(GameMap.Fields[x/FIELD_SIZE]))*FIELD_SIZE {
		return true
	}
	arrayPosX := x / FIELD_SIZE
	arrayPosY := y / FIELD_SIZE
	accessible0, accessible1 := true, true
	if GameMap.Fields[arrayPosX][arrayPosY].Contains[0] != nil {
		if GameMap.Fields[arrayPosX][arrayPosY].Contains[0].getType() == 1 || GameMap.Fields[arrayPosX][arrayPosY].Contains[0].getType() == 10 || GameMap.Fields[arrayPosX][arrayPosY].Contains[0].getType() == 11 {
			return true
		}
		accessible0 = GameMap.Fields[arrayPosX][arrayPosY].Contains[0].isAccessible()
	}
	if GameMap.Fields[arrayPosX][arrayPosY].Contains[1] != nil {
		if GameMap.Fields[arrayPosX][arrayPosY].Contains[1].getType() == 1 || GameMap.Fields[arrayPosX][arrayPosY].Contains[1].getType() == 10 || GameMap.Fields[arrayPosX][arrayPosY].Contains[1].getType() == 11 {
			return true
		}
		accessible1 = GameMap.Fields[arrayPosX][arrayPosY].Contains[1].isAccessible()
	}
	isAccessible := accessible0 && accessible1
	return isAccessible
}

/*
???
*/
func (b *Bomberman) collisionWithSurroundings(xOffset int, yOffset int) bool {
	topRight := outerEdges(b.topRightPos.x+xOffset, b.topRightPos.y+yOffset)
	topLeft := outerEdges(b.topLeftPos.x+xOffset, b.topLeftPos.y+yOffset)
	bottomRight := outerEdges(b.bottomRightPos.x+xOffset, b.bottomRightPos.y+yOffset)
	bottomLeft := outerEdges(b.bottomLeftPos.x+xOffset, b.bottomLeftPos.y+yOffset)
	legal := topRight && topLeft && bottomRight && bottomLeft
	//check a start  of method?
	if b.GhostActive {
		return true
	} else {
		return legal
	}
}
