package main

import (
	"container/list"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
	"time"
)

//commit comment
const FIELD_SIZE = 50
const STEP_SIZE = 3
const CANVAS_SIZE = 1000

var GameMap = NewMap(CANVAS_SIZE / FIELD_SIZE)
var Connections = make(map[uint64]*Session, 0)
var ticker = time.NewTicker(16 * time.Millisecond)
var incomingTicker = time.NewTicker(1 * time.Millisecond)

//Things send to the clients
var bombermanArray = make([]Bomberman, 0)
var abstractGameMap = make([][][]FieldObject, 0)
var clientPackageAsJson = make([]byte, 0)

type Position struct {
	x int
	y int
}

func newPosition(x int, y int) Position {
	return Position{
		x: x,
		y: y,
	}
}

type KeyInput struct {
	Wpressed     bool `json:"w"`
	Spressed     bool `json:"s"`
	Apressed     bool `json:"a"`
	Dpressed     bool `json:"d"`
	SpacePressed bool `json:" "`
}

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

type ClientPackage struct {
	Players []Bomberman
	GameMap [][][]FieldObject
	//TestPlayer [][]int
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

//Wrapper for the user
type Session struct {
	User              *User           //Connected user
	Bomber            *Bomberman      //Bomber of the connected user
	Connection        *websocket.Conn //Websocket connection
	ConnectionStarted time.Time       //point when player joined
}

func NewSession(user *User, character *Bomberman, connection *websocket.Conn, connectionStarted time.Time) *Session {
	return &Session{User: user, Bomber: character, Connection: connection, ConnectionStarted: connectionStarted}
}

//Returns the string representation of the connection
func (r *Session) String() string {
	return "Session: { " + r.User.String() + "|" + r.Bomber.String() + "|" + r.Connection.RemoteAddr().String() + "|" + r.ConnectionStarted.String() + "}"
}

//Prints every active connection
func AllConnectionsAsString() string {
	result := "Active Connections:"
	for _, v := range Connections {
		result += v.String() + "\n"
	}
	return result
}

//Starts the interaction loop
func StartPlayerLoop(session *Session) {
	//Add the infos to the connection map
	BuildAbstractGameMap()
	Connections[session.User.UserID] = session
	GameMap.Fields[(session.Bomber.PositionX+FIELD_SIZE/2)/FIELD_SIZE][(session.Bomber.PositionY+FIELD_SIZE/2)/FIELD_SIZE].Player.PushBack(session.Bomber)
	playerWebsocketLoop(session)
	//Remove player from list at his last array position
	x := (session.Bomber.PositionX + FIELD_SIZE/2) / FIELD_SIZE
	y := (session.Bomber.PositionY + FIELD_SIZE/2) / FIELD_SIZE
	removePlayerFromList(GameMap.Fields[x][y].Player, session.Bomber)
	//Remove from the connection map
	delete(Connections, session.User.UserID)
}

//interaction loop
func playerWebsocketLoop(session *Session) {
	for {
		_, p, err := session.Connection.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		var keys KeyInput
		if err := json.Unmarshal(p, &keys); err != nil {
			log.Println(err)
			continue
		}
		session.Bomber.IsMoving = false
		if keys.Wpressed {
			session.Bomber.DirUp, session.Bomber.DirDown, session.Bomber.DirLeft, session.Bomber.DirRight = true, false, false, false
			session.Bomber.IsMoving = true
			if session.Bomber.collisionWithSurroundings(0, -int(STEP_SIZE*session.Bomber.stepMult)) {
				if session.Bomber.isMovementLegal(session.Bomber.PositionX, session.Bomber.PositionY-int(STEP_SIZE*session.Bomber.stepMult)) {

					session.Bomber.topRightPos.updatePosition(0, -int(STEP_SIZE*session.Bomber.stepMult))
					session.Bomber.topLeftPos.updatePosition(0, -int(STEP_SIZE*session.Bomber.stepMult))
					session.Bomber.bottomRightPos.updatePosition(0, -int(STEP_SIZE*session.Bomber.stepMult))
					session.Bomber.bottomLeftPos.updatePosition(0, -int(STEP_SIZE*session.Bomber.stepMult))
					session.Bomber.PositionY -= int(STEP_SIZE * session.Bomber.stepMult)
				}
			}
		} else
		//S
		if keys.Spressed {
			session.Bomber.DirUp, session.Bomber.DirDown, session.Bomber.DirLeft, session.Bomber.DirRight = false, true, false, false
			session.Bomber.IsMoving = true
			if session.Bomber.collisionWithSurroundings(0, int(STEP_SIZE*session.Bomber.stepMult)) {
				if session.Bomber.isMovementLegal(session.Bomber.PositionX, session.Bomber.PositionY+int(STEP_SIZE*session.Bomber.stepMult)) {

					session.Bomber.topRightPos.updatePosition(0, int(STEP_SIZE*session.Bomber.stepMult))
					session.Bomber.topLeftPos.updatePosition(0, int(STEP_SIZE*session.Bomber.stepMult))
					session.Bomber.bottomRightPos.updatePosition(0, int(STEP_SIZE*session.Bomber.stepMult))
					session.Bomber.bottomLeftPos.updatePosition(0, int(STEP_SIZE*session.Bomber.stepMult))
					session.Bomber.PositionY += int(STEP_SIZE * session.Bomber.stepMult)
				}
			}
		} else
		//A
		if keys.Apressed {
			session.Bomber.DirUp, session.Bomber.DirDown, session.Bomber.DirLeft, session.Bomber.DirRight = false, false, true, false
			session.Bomber.IsMoving = true
			if session.Bomber.collisionWithSurroundings(-int(STEP_SIZE*session.Bomber.stepMult), 0) {
				if session.Bomber.isMovementLegal(session.Bomber.PositionX-int(STEP_SIZE*session.Bomber.stepMult), session.Bomber.PositionY) {

					session.Bomber.topRightPos.updatePosition(-int(STEP_SIZE*session.Bomber.stepMult), 0)
					session.Bomber.topLeftPos.updatePosition(-int(STEP_SIZE*session.Bomber.stepMult), 0)
					session.Bomber.bottomRightPos.updatePosition(-int(STEP_SIZE*session.Bomber.stepMult), 0)
					session.Bomber.bottomLeftPos.updatePosition(-int(STEP_SIZE*session.Bomber.stepMult), 0)
					session.Bomber.PositionX -= int(STEP_SIZE * session.Bomber.stepMult)
				}
			}
		} else
		//D
		if keys.Dpressed {
			session.Bomber.DirUp, session.Bomber.DirDown, session.Bomber.DirLeft, session.Bomber.DirRight = false, false, false, true
			session.Bomber.IsMoving = true
			if session.Bomber.collisionWithSurroundings(int(STEP_SIZE*session.Bomber.stepMult), 0) {
				if session.Bomber.isMovementLegal(session.Bomber.PositionX+int(STEP_SIZE*session.Bomber.stepMult), session.Bomber.PositionY) {

					session.Bomber.topRightPos.updatePosition(int(STEP_SIZE*session.Bomber.stepMult), 0)
					session.Bomber.topLeftPos.updatePosition(int(STEP_SIZE*session.Bomber.stepMult), 0)
					session.Bomber.bottomRightPos.updatePosition(int(STEP_SIZE*session.Bomber.stepMult), 0)
					session.Bomber.bottomLeftPos.updatePosition(int(STEP_SIZE*session.Bomber.stepMult), 0)
					session.Bomber.PositionX += int(STEP_SIZE * session.Bomber.stepMult)
				}
			}
		}
		//Spacebar
		if keys.SpacePressed {
			if session.Bomber.IsAlive {
				go session.Bomber.placeBomb()
			}
		}
		//if session.Bomber.IsAlive && !itemActive {
		//	checkItem(session)
		//}

	}

}

//func checkItem(session *Session) {
//	arrayPosX := (session.Bomber.PositionX + FIELD_SIZE/2) / FIELD_SIZE
//	arrayPosY := (session.Bomber.PositionY + FIELD_SIZE/2) / FIELD_SIZE
//
//	//!time.AfterFunc calls new goroutine automatically!
//	for i := 0; i < len(GameMap.Fields[arrayPosX][arrayPosY].Contains); i++ {
//		if GameMap.Fields[arrayPosX][arrayPosY].Contains[i] != nil {
//			//
//			if GameMap.Fields[arrayPosX][arrayPosY].Contains[i].getType() == 6 {
//				itemActive = true
//				session.Bomber.stepMult = 1.5
//				GameMap.Fields[arrayPosX][arrayPosY].Contains[i] = nil
//				BuildAbstractGameMap()
//				time.AfterFunc(7*time.Second, func() {
//					if session.Bomber.IsAlive {
//						session.Bomber.stepMult = 1
//					}
//					itemActive = false
//				})
//			} else if GameMap.Fields[arrayPosX][arrayPosY].Contains[i].getType() == 7 {
//				itemActive = true
//				session.Bomber.stepMult = 0.5
//				GameMap.Fields[arrayPosX][arrayPosY].Contains[i] = nil
//				BuildAbstractGameMap()
//				time.AfterFunc(5*time.Second, func() {
//					if session.Bomber.IsAlive {
//						session.Bomber.stepMult = 1
//					}
//					itemActive = false
//				})
//
//			} else if GameMap.Fields[arrayPosX][arrayPosY].Contains[i].getType() == 8 {
//				itemActive = true
//				session.Bomber.GhostActive = true
//				GameMap.Fields[arrayPosX][arrayPosY].Contains[i] = nil
//				BuildAbstractGameMap()
//				time.AfterFunc(5*time.Second, func() {
//					if session.Bomber.IsAlive {
//						session.Bomber.GhostActive = false
//					}
//					itemActive = false
//				})
//			} else if GameMap.Fields[arrayPosX][arrayPosY].Contains[i].getType() == 12 {
//				removePlayerFromList(GameMap.Fields[arrayPosX][arrayPosY].Player, session.Bomber)
//				session.Bomber.PositionX = 948
//				session.Bomber.PositionY = 948
//				GameMap.Fields[arrayPosX][arrayPosY].Player.PushBack(session.Bomber)
//				BuildAbstractGameMap()
//			}
//		}
//	}
//}

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
	case 12:

	default:
		return

	}
	BuildAbstractGameMap()
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
		BuildAbstractGameMap()
		r.hasTeleported = true
	}
}

func (r *Bomberman) isMovementLegal(x int, y int) bool {
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
				removePlayerFromList(GameMap.Fields[oldPosX][arrayPosY].Player, r)
				GameMap.Fields[arrayPosX][arrayPosY].Player.PushBack(r)
				r.checkFieldForItem(arrayPosX, arrayPosY)
				if r.hasTeleported {
					r.hasTeleported = false
					return false
				}
				return true
			} else {
				return false
			}
		} else if oldPosY != arrayPosY {
			if r.isFieldAccessible(x, y) || r.GhostActive {
				removePlayerFromList(GameMap.Fields[arrayPosX][oldPosY].Player, r)
				GameMap.Fields[arrayPosX][arrayPosY].Player.PushBack(r)
				r.checkFieldForItem(arrayPosX, arrayPosY)
				if r.hasTeleported {
					r.hasTeleported = false
					return false
				}
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

func outerEdges(x int, y int) bool {
	if x < 0 || y < 0 || x > (len(GameMap.Fields)-1)*FIELD_SIZE || y > (len(GameMap.Fields[x/FIELD_SIZE])-1)*FIELD_SIZE {
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

func (b *Bomberman) collisionWithSurroundings(xOffset int, yOffset int) bool {
	topRight := outerEdges(b.topRightPos.x+xOffset, b.topRightPos.y+yOffset)
	topLeft := outerEdges(b.topLeftPos.x+xOffset, b.topLeftPos.y+yOffset)
	bottomRight := outerEdges(b.bottomRightPos.x+xOffset, b.bottomRightPos.y+yOffset)
	bottomLeft := outerEdges(b.bottomLeftPos.x+xOffset, b.bottomLeftPos.y+yOffset)
	legal := topRight && topLeft && bottomRight && bottomLeft
	if b.GhostActive {
		return true
	} else {
		return legal
	}
}

func printList(list *list.List) {
	element := list.Front()
	if element == nil {
		log.Println("List is null!")
		return
	}
	log.Println("List started: ")
	log.Println(element.Value.(*Bomberman))
	for element.Next() != nil {
		log.Println(element.Value.(*Bomberman))
		element = element.Next()
	}
	log.Println("List ended...")
}

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

func UpdateClients() {
	for _ = range ticker.C {
		err := sendDataToClients()
		if err != nil {
			log.Println(err)
			break
		}
	}
	log.Println("Updating Clients stopped.")
}

func sendDataToClients() error {
	//Create array from all connected Bombermen
	bombermanArray = make([]Bomberman, len(Connections))
	count := 0

	for _, v := range Connections {
		bombermanArray[count] = *v.Bomber
		count++
	}

	var err error
	clientPackageAsJson, err = json.Marshal(ClientPackage{
		Players: bombermanArray,
		GameMap: abstractGameMap,
		//TestPlayer: testToSend,
	})
	if err != nil {
		log.Println(err)
		return err
	}

	for _, v := range Connections {
		if err := v.Connection.WriteMessage(websocket.TextMessage, clientPackageAsJson); err != nil {
			return err
		}
	}
	return nil
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

func isLesserThan(a interface{}, b interface{}) bool {
	return a.(*Session).User.UserID < b.(*Session).User.UserID
}

func (p *Position) updatePosition(xOffset int, yOffset int) {
	p.x += xOffset
	p.y += yOffset
}
