package main

import (
	"github.com/gorilla/websocket"
	"log"
	"time"
)

var connections = make(map[uint64]*Session, 0)

//Wrapper for the user
type Session struct {
	User              *User           //Connected user
	Character         *Character      //Character of the connected user
	Connection        *websocket.Conn //Websocket connection
	ConnectionStarted time.Time       //point when player joined
}

func NewSession(user *User, character *Character, connection *websocket.Conn, connectionStarted time.Time) *Session {
	return &Session{User: user, Character: character, Connection: connection, ConnectionStarted: connectionStarted}
}

//Returns the string representation of the connection
func (r *Session) String() string {
	return "Session: { " + r.User.String() + "|" + r.Connection.RemoteAddr().String() + "|" + r.ConnectionStarted.String() + "}"
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
		messageType, p, err := session.Connection.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(p)
		if err := session.Connection.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

func sendNewPositionToClients() {

}
