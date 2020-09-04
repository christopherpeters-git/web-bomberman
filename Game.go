package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
	"time"
)

var connections = make(map[uint64]*Session, 0)

var ticker = time.NewTicker(5 * time.Millisecond)

type Bomberman struct {
	UserID         uint64
	PositionX      int
	PositionY      int
	lastBombPlaced time.Time
}

func (r *Bomberman) String() string {
	return "Bomberman: {" + strconv.FormatUint(r.UserID, 10) + " | " + strconv.FormatInt(int64(r.PositionX), 10) + " | " + strconv.FormatInt(int64(r.PositionY), 10) + " | " + r.lastBombPlaced.String() + "}"
}

func NewBomberman(userID uint64, positionX int, positionY int) *Bomberman {
	return &Bomberman{UserID: userID, PositionX: positionX, PositionY: positionY}
}

//Wrapper for the user
type Session struct {
	User              *User           //Connected user
	Character         *Bomberman      //Character of the connected user
	Connection        *websocket.Conn //Websocket connection
	ConnectionStarted time.Time       //point when player joined
}

func NewSession(user *User, character *Bomberman, connection *websocket.Conn, connectionStarted time.Time) *Session {
	return &Session{User: user, Character: character, Connection: connection, ConnectionStarted: connectionStarted}
}

//Returns the string representation of the connection
func (r *Session) String() string {
	return "Session: { " + r.User.String() + "|" + r.Character.String() + "|" + r.Connection.RemoteAddr().String() + "|" + r.ConnectionStarted.String() + "}"
}

//Prints every active connection
func AllConnectionsAsString() string {
	result := "Active Connections:"
	for _, v := range connections {
		result += v.String() + "\n"
	}
	return result
}

//Starts the interaction loop
func StartPlayerLoop(session *Session) {
	//Add the infos to the connection map
	connections[session.User.UserID] = session
	playerWebsocketLoop(session)
	//Remove from the connection map
	delete(connections, session.User.UserID)
}

//interaction loop
func playerWebsocketLoop(session *Session) {
	for {
		_, p, err := session.Connection.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("incoming (unformatted): " + string(p))
		switch string(p) {
		//W
		case "w":
			session.Character.PositionY += 10
			break
		//A
		case "a":
			session.Character.PositionX -= 10
			break
		//S
		case "s":
			session.Character.PositionY -= 10
			break
		//D
		case "d":
			session.Character.PositionX += 10
			break
		default:
			break
		}
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
	//collect data
	sessions := make([]Bomberman, len(connections))
	count := 0
	for _, v := range connections {
		sessions[count] = *v.Character
		count++
	}
	jsonBytes, err := json.MarshalIndent(sessions, "", " ")
	if err != nil {
		return err
	}
	for _, v := range connections {
		if err := v.Connection.WriteMessage(websocket.TextMessage, jsonBytes); err != nil {
			return err
		}
	}
	return nil
}
