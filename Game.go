package main

import (
	"container/list"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

//commit comment
const (
	FIELD_SIZE                  = 50
	STEP_SIZE                   = 3
	CANVAS_SIZE                 = 1000
	STANDARD_BOMB_RADIUS        = 2
	STANDARD_BOMB_TIME          = 3
	STANDARD_STEP_MULTIPLICATOR = 1
	DEATH_STEP_MULTIPLICATOR    = 0.5
)

var GameMap = NewMap(CANVAS_SIZE / FIELD_SIZE)
var Connections = make(map[uint64]*Session, 0)
var ticker = time.NewTicker(16 * time.Millisecond)
var spawnPositions = [][]int{{0, 0}, {0, 10}, {0, 19}, {10, 0}, {10, 19}, {19, 0}, {19, 10}, {19, 19}}

//var incomingTicker = time.NewTicker(1 * time.Millisecond)
var sessionRunning = false

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

func pixToArr(pixel int) int {
	return (pixel + FIELD_SIZE/2) / FIELD_SIZE
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

type KeyInput struct {
	Wpressed     bool `json:"w"`
	Spressed     bool `json:"s"`
	Apressed     bool `json:"a"`
	Dpressed     bool `json:"d"`
	SpacePressed bool `json:" "`
}

type ClientPackage struct {
	Players []Bomberman
	GameMap [][][]FieldObject
	//TestPlayer [][]int
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
		realStepSize := STEP_SIZE * session.Bomber.stepMult
		if keys.Wpressed {
			session.Bomber.DirUp, session.Bomber.DirDown, session.Bomber.DirLeft, session.Bomber.DirRight = true, false, false, false
			session.Bomber.IsMoving = true
			if session.Bomber.collisionWithSurroundings(0, -int(realStepSize)) {
				if session.Bomber.isMovementLegal(session.Bomber.PositionX, session.Bomber.PositionY-int(realStepSize)) {

					session.Bomber.topRightPos.updatePosition(0, -int(realStepSize))
					session.Bomber.topLeftPos.updatePosition(0, -int(realStepSize))
					session.Bomber.bottomRightPos.updatePosition(0, -int(realStepSize))
					session.Bomber.bottomLeftPos.updatePosition(0, -int(realStepSize))
					session.Bomber.PositionY -= int(realStepSize)
				}
			}
		} else
		//S
		if keys.Spressed {
			session.Bomber.DirUp, session.Bomber.DirDown, session.Bomber.DirLeft, session.Bomber.DirRight = false, true, false, false
			session.Bomber.IsMoving = true
			if session.Bomber.collisionWithSurroundings(0, int(realStepSize)) {
				if session.Bomber.isMovementLegal(session.Bomber.PositionX, session.Bomber.PositionY+int(realStepSize)) {

					session.Bomber.topRightPos.updatePosition(0, int(realStepSize))
					session.Bomber.topLeftPos.updatePosition(0, int(realStepSize))
					session.Bomber.bottomRightPos.updatePosition(0, int(realStepSize))
					session.Bomber.bottomLeftPos.updatePosition(0, int(realStepSize))
					session.Bomber.PositionY += int(realStepSize)
				}
			}
		} else
		//A
		if keys.Apressed {
			session.Bomber.DirUp, session.Bomber.DirDown, session.Bomber.DirLeft, session.Bomber.DirRight = false, false, true, false
			session.Bomber.IsMoving = true
			if session.Bomber.collisionWithSurroundings(-int(realStepSize), 0) {
				if session.Bomber.isMovementLegal(session.Bomber.PositionX-int(realStepSize), session.Bomber.PositionY) {

					session.Bomber.topRightPos.updatePosition(-int(realStepSize), 0)
					session.Bomber.topLeftPos.updatePosition(-int(realStepSize), 0)
					session.Bomber.bottomRightPos.updatePosition(-int(realStepSize), 0)
					session.Bomber.bottomLeftPos.updatePosition(-int(realStepSize), 0)
					session.Bomber.PositionX -= int(realStepSize)
				}
			}
		} else
		//D
		if keys.Dpressed {
			session.Bomber.DirUp, session.Bomber.DirDown, session.Bomber.DirLeft, session.Bomber.DirRight = false, false, false, true
			session.Bomber.IsMoving = true
			if session.Bomber.collisionWithSurroundings(int(realStepSize), 0) {
				if session.Bomber.isMovementLegal(session.Bomber.PositionX+int(realStepSize), session.Bomber.PositionY) {

					session.Bomber.topRightPos.updatePosition(int(realStepSize), 0)
					session.Bomber.topLeftPos.updatePosition(int(realStepSize), 0)
					session.Bomber.bottomRightPos.updatePosition(int(realStepSize), 0)
					session.Bomber.bottomLeftPos.updatePosition(int(realStepSize), 0)
					session.Bomber.PositionX += int(realStepSize)
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

func isLesserThan(a interface{}, b interface{}) bool {
	return a.(*Session).User.UserID < b.(*Session).User.UserID
}

func (p *Position) updatePosition(xOffset int, yOffset int) {
	p.x += xOffset
	p.y += yOffset
}

func StartGameIfPlayersReady() {
	if len(Connections) < 2 {
		return
	}
	for _, v := range Connections {
		if !v.Bomber.PlayerReady {
			return
		}
	}
	resetGame("images/map2.png")
	sessionRunning = true
	for _, v := range Connections {
		v.Bomber.PlayerReady = false
	}

}

func resetGame(s string) {
	playerDied = false
	GameMap.clear()
	CreateMapFromImage(GameMap, s)
	count := 0
	for _, v := range Connections {
		if count > 7 {
			count = 0
		}
		v.Bomber.Reset(spawnPositions[count][0], spawnPositions[count][1])
		count++
	}
}

func killAllPlayersOnField(list *list.List) {
	element := list.Front()
	if element != nil {
		element.Value.(*Bomberman).Kill()
		playerDied = true
		for element.Next() != nil {
			element = element.Next()
			element.Value.(*Bomberman).Kill()
		}
	}
}

func isOnePlayerAlive() {
	counter := 0
	var lastBomberAlive *Bomberman
	for _, v := range Connections {
		if v.Bomber.IsAlive {
			lastBomberAlive = v.Bomber
			counter++
		}
	}
	if counter > 1 {
		return
	} else if counter == 0 {
		log.Println("Draw")
	} else if counter == 1 {
		log.Println(lastBomberAlive.Name)
		log.Println("has Won")
	}
	//todo send message
	resetGame("images/map.png")
	sessionRunning = false
}
