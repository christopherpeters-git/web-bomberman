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
var itemActive = false

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
	GhostActive    bool
	topRightPos    Position
	topLeftPos     Position
	bottomRightPos Position
	bottomLeftPos  Position
	stepMult       float32
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
		BombRadius:     3,
		bombTime:       3,
		IsAlive:        true,
		IsHit:          false,
		GhostActive:    false,
		topRightPos:    newPosition(43, 7),
		topLeftPos:     newPosition(7, 7),
		bottomRightPos: newPosition(7, 43),
		bottomLeftPos:  newPosition(43, 43),
		stepMult:       1,
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
	GameMap.Fields[0][0].Player.PushBack(session.Bomber)
	playerWebsocketLoop(session)
	//Remove player from list at his last array position
	x := session.Bomber.PositionX / FIELD_SIZE
	y := session.Bomber.PositionY / FIELD_SIZE
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
		//if !session.Bomber.IsAlive {
		//	return
		//}
		if keys.Wpressed {
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
		if session.Bomber.IsAlive && !itemActive {
			checkItem(session)
		}
	}

}

func checkItem(session *Session) {
	arrayPosX := (session.Bomber.PositionX + FIELD_SIZE/2) / FIELD_SIZE
	arrayPosY := (session.Bomber.PositionY + FIELD_SIZE/2) / FIELD_SIZE

	//!time.AfterFunc calls new goroutine automatically!
	for i := 0; i < len(GameMap.Fields[arrayPosX][arrayPosY].Contains); i++ {
		if GameMap.Fields[arrayPosX][arrayPosY].Contains[i] != nil {
			if GameMap.Fields[arrayPosX][arrayPosY].Contains[i].getType() == 6 {
				//resets set timer, so the duration of item collected before doesnt count anymore
				itemActive = true
				session.Bomber.stepMult = 1.8
				GameMap.Fields[arrayPosX][arrayPosY].Contains[i] = nil
				BuildAbstractGameMap()
				time.AfterFunc(7*time.Second, func() {

					itemActive = false
				})
			} else if GameMap.Fields[arrayPosX][arrayPosY].Contains[i].getType() == 7 {
				itemActive = true
				session.Bomber.stepMult = 0.5
				GameMap.Fields[arrayPosX][arrayPosY].Contains[i] = nil
				BuildAbstractGameMap()
				time.AfterFunc(5*time.Second, func() {
					if session.Bomber.IsAlive {
						session.Bomber.stepMult = 1
					}
					itemActive = false
				})

			} else if GameMap.Fields[arrayPosX][arrayPosY].Contains[i].getType() == 8 {
				itemActive = true
				session.Bomber.GhostActive = true
				GameMap.Fields[arrayPosX][arrayPosY].Contains[i] = nil
				BuildAbstractGameMap()
				time.AfterFunc(5*time.Second, func() {
					if session.Bomber.IsAlive {
						session.Bomber.GhostActive = false
					}
					itemActive = false
				})

			}

		}
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
				return true
			} else {
				return false
			}
		} else if oldPosY != arrayPosY {
			if r.isFieldAccessible(x, y) || r.GhostActive {
				removePlayerFromList(GameMap.Fields[arrayPosX][oldPosY].Player, r)
				GameMap.Fields[arrayPosX][arrayPosY].Player.PushBack(r)
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
	if isAccessible {
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
		if GameMap.Fields[arrayPosX][arrayPosY].Contains[0].getType() == 1 {
			return true
		}
		accessible0 = GameMap.Fields[arrayPosX][arrayPosY].Contains[0].isAccessible()
	}
	if GameMap.Fields[arrayPosX][arrayPosY].Contains[1] != nil {
		if GameMap.Fields[arrayPosX][arrayPosY].Contains[1].getType() == 1 {
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
