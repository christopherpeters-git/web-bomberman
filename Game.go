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
	for {
		sendDataToClients()
	}
}

func sendDataToClients() {
	//collect data
	sessions := make([]Session, len(connections))
	count := 0
	for _, v := range connections {
		sessions[count] = *v
		count++
	}
	//jsonBytes, err := json.Marshal(sessions)
	//send data to all clients
	//for _, v := range connections {
	//	if err := v.Connection.WriteMessage(websocket.TextMessage, p); err != nil {
	//		log.Println(err)
	//		return
	//	}
	//}
}
