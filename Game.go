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
const STEP_SIZE = 4
const CANVAS_SIZE = 500

var GameMap = NewMap(CANVAS_SIZE / FIELD_SIZE)
var connections = sortedmap.New(10, isLesserThan)
var ticker = time.NewTicker(5 * time.Millisecond)

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
	OldPositionX   int
	oldPositionY   int
	lastBombPlaced time.Time
	BombRadius     int
	bombTime       int
	IsAlive        bool
}

type ClientPackage struct {
	Players    []Bomberman
	GameMap    [][][]FieldObject
	TestPlayer [][]int
}

func (r *Bomberman) String() string {
	return "Bomberman: {" + strconv.FormatUint(r.UserID, 10) + " | " + strconv.FormatInt(int64(r.PositionX), 10) + " | " + strconv.FormatInt(int64(r.PositionY), 10) + " | " + r.lastBombPlaced.String() + "}"
}

func NewBomberman(userID uint64, positionX int, positionY int, name string) *Bomberman {
	return &Bomberman{
		UserID:       userID,
		PositionX:    positionX,
		PositionY:    positionY,
		OldPositionX: positionX,
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
		var keys KeyInput
		if err := json.Unmarshal(p, &keys); err != nil {
			log.Println(err)
			continue
		}
		//if !session.Bomber.IsAlive {
		//	return
		//}
		if keys.Wpressed {
			if session.Bomber.isMovementLegal(session.Bomber.PositionX, session.Bomber.PositionY-STEP_SIZE) {

				session.Bomber.PositionY -= STEP_SIZE

			}
		} else
		//S
		if keys.Spressed {
			if session.Bomber.isMovementLegal(session.Bomber.PositionX, session.Bomber.PositionY+STEP_SIZE) {

				session.Bomber.PositionY += STEP_SIZE

			}
		} else
		//A
		if keys.Apressed {
			if session.Bomber.isMovementLegal(session.Bomber.PositionX-STEP_SIZE, session.Bomber.PositionY) {

				session.Bomber.PositionX -= STEP_SIZE

			}
		} else
		//D
		if keys.Dpressed {
			if session.Bomber.isMovementLegal(session.Bomber.PositionX+STEP_SIZE, session.Bomber.PositionY) {
				session.Bomber.PositionX += STEP_SIZE

			}
		}
		//Spacebar
		if keys.SpacePressed {
			go session.Bomber.placeBomb()
		}
	}

}

func (r *Bomberman) isMovementLegal(x int, y int) bool { //r.positionX = 50
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
			if r.isFieldAccessible(x, y) {
				removePlayerFromList(GameMap.Fields[oldPosX][arrayPosY].Player, r)
				GameMap.Fields[arrayPosX][arrayPosY].Player.PushBack(r)
				return true
			} else {
				return false
			}
		} else if oldPosY != arrayPosY {
			if r.isFieldAccessible(x, y) {
				removePlayerFromList(GameMap.Fields[arrayPosX][oldPosY].Player, r)
				GameMap.Fields[arrayPosX][arrayPosY].Player.PushBack(r)
				return true
			} else {
				return false
			}
		}
		r.OldPositionX = r.PositionX
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
		b.OldPositionX = b.PositionX
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
	//Create array from all connected Bombermen
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
	//Create map to send
	mapToSend := make([][][]FieldObject, len(GameMap.Fields))
	testToSend := make([][]int, len(GameMap.Fields))
	for i, _ := range GameMap.Fields {
		mapToSend[i] = make([][]FieldObject, len(GameMap.Fields[i]))
		testToSend[i] = make([]int, len(GameMap.Fields[i]))
		for j, _ := range GameMap.Fields[i] {
			mapToSend[i][j] = make([]FieldObject, len(GameMap.Fields[i][j].Contains))
			if GameMap.Fields[i][j].Player.Front() != nil {
				testToSend[i][j] = 1
			}
			for k, _ := range GameMap.Fields[i][j].Contains {
				if GameMap.Fields[i][j].Contains[k] != nil {
					mapToSend[i][j][k] = GameMap.Fields[i][j].Contains[k].getType()
				}
			}
		}
	}

	//Create ClientPackage to send to every client
	clientPackage := ClientPackage{
		Players:    sessions,
		GameMap:    mapToSend,
		TestPlayer: testToSend,
	}

	jsonBytes, err := json.MarshalIndent(clientPackage, "", " ")
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
