package main

import (
	"container/list"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/umpc/go-sortedmap"
	"log"
	"strconv"
	"time"
)

//commit comment
const FIELD_SIZE = 50
const STEP_SIZE = 10
const CANVAS_SIZE = 500

var GameMap = NewMap(CANVAS_SIZE / FIELD_SIZE)
var connections = sortedmap.New(10, isLesserThan)
var ticker = time.NewTicker(5 * time.Millisecond)

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
}

func (r *Bomberman) String() string {
	return "Bomberman: {" + strconv.FormatUint(r.UserID, 10) + " | " + strconv.FormatInt(int64(r.PositionX), 10) + " | " + strconv.FormatInt(int64(r.PositionY), 10) + " | " + r.lastBombPlaced.String() + "}"
}

func NewBomberman(userID uint64, positionX int, positionY int, name string) *Bomberman {
	return &Bomberman{
		UserID:       userID,
		PositionX:    positionX,
		PositionY:    positionY,
		oldPositionX: positionX,
		oldPositionY: positionY,
		Name:         name,
		BombRadius:   3,
		bombTime:     3,
		IsAlive:      true,
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

	iterCh, err := connections.IterCh()

	if err != nil {
		log.Println(err)
		return result
	}
	defer iterCh.Close()

	for v := range iterCh.Records() {
		result += v.Val.(*Session).String() + "\n"
	}
	return result
}

//Starts the interaction loop
func StartPlayerLoop(session *Session) {
	//Add the infos to the connection map
	connections.Insert(session.User.UserID, session)
	//FillTestMap(GameMap)
	GameMap.Fields[0][0].Player.PushBack(session.Bomber)
	playerWebsocketLoop(session)
	//Remove from the connection map
	connections.Delete(session.User.UserID)
}

//interaction loop
func playerWebsocketLoop(session *Session) {
	for {
		_, p, err := session.Connection.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		if !session.Bomber.IsAlive {
			return
		}

		switch string(p) {
		//W
		case "w":
			if session.Bomber.canEnter(session.Bomber.PositionX, session.Bomber.PositionY-STEP_SIZE) {
				session.Bomber.PositionY -= STEP_SIZE
				updatePlayerPositioning(session)
			}

		//A
		case "a":
			if session.Bomber.canEnter(session.Bomber.PositionX-STEP_SIZE, session.Bomber.PositionY) {
				session.Bomber.PositionX -= STEP_SIZE
				updatePlayerPositioning(session)
			}

		//S
		case "s":
			if session.Bomber.canEnter(session.Bomber.PositionX, session.Bomber.PositionY+STEP_SIZE) {
				session.Bomber.PositionY += STEP_SIZE
				updatePlayerPositioning(session)
			}

		//D
		case "d":
			if session.Bomber.canEnter(session.Bomber.PositionX+STEP_SIZE, session.Bomber.PositionY) {
				session.Bomber.PositionX += STEP_SIZE
				updatePlayerPositioning(session)
			}
		//Spacebar
		case " ":
			go session.Bomber.placeBomb()

		default:
			break
		}

	}

}
func updatePlayerPositioning(session *Session) {
	posX := session.Bomber.PositionX / FIELD_SIZE
	posY := session.Bomber.PositionY / FIELD_SIZE
	oldPosX := session.Bomber.oldPositionX / FIELD_SIZE
	oldPosY := session.Bomber.oldPositionY / FIELD_SIZE
	//Change Pushback
	if posX != oldPosX {
		if session.Bomber.isFieldAccessible() {
			removePlayerFromList(GameMap.Fields[oldPosX][posY].Player, session.Bomber)
			GameMap.Fields[posX][posY].Player.PushBack(session.Bomber)
			//log.Println(GameMap.Fields[posX][posY].Player)
		}
	} else if posY != oldPosY {
		if session.Bomber.isFieldAccessible() {
			removePlayerFromList(GameMap.Fields[posX][oldPosY].Player, session.Bomber)
			GameMap.Fields[posX][posY].Player.PushBack(session.Bomber)
			//log.Println(GameMap.Fields[posX][posY].Player)
		}
	}

}

func (r *Bomberman) canEnter(x int, y int) bool {
	if x < 0 || y < 0 || x >= len(GameMap.Fields)*FIELD_SIZE || y >= len(GameMap.Fields[x/FIELD_SIZE])*FIELD_SIZE {
		return false
	}
	arrayPosX := x / FIELD_SIZE
	arrayPosY := y / FIELD_SIZE
	inBounds := arrayPosX >= 0 && arrayPosY >= 0 && arrayPosX < len(GameMap.Fields) && arrayPosY < len(GameMap.Fields[arrayPosX])
	return inBounds
}

func (b *Bomberman) isFieldAccessible() bool {
	isAccessNull := true
	isAccessOne := true
	arrayPosX := b.PositionX / FIELD_SIZE
	arrayPosY := b.PositionY / FIELD_SIZE
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
	//collect data
	sessions := make([]Bomberman, connections.Len())
	count := 0

	iterCh, err := connections.IterCh()

	if err != nil {
		return nil
	}
	defer iterCh.Close()

	for v := range iterCh.Records() {
		sessions[count] = *v.Val.(*Session).Bomber
		count++
	}

	jsonBytes, err := json.MarshalIndent(sessions, "", " ")
	if err != nil {

		return err
	}
	iterCh, err = connections.IterCh()

	if err != nil {
		return nil
	}

	for v := range iterCh.Records() {

		if err := v.Val.(*Session).Connection.WriteMessage(websocket.TextMessage, jsonBytes); err != nil {
			return err
		}
	}
	return nil
}

func isLesserThan(a interface{}, b interface{}) bool {
	return a.(*Session).User.UserID < b.(*Session).User.UserID
}
