package main

import (
	global "./global"
	"database/sql"
	"github.com/gorilla/websocket"
	_ "github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
	//Error:             nil,
}

func StartWebSocketConnection(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	//Check if db connection is available
	if err := db.Ping(); err != nil {
		http.Error(w, global.INTERNAL_SERVER_ERROR_RESPONSE, http.StatusInternalServerError)
		log.Println(err)
		return
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	//TODO read username and password?
	user, dErr := GetUserFromDB(db, "TEST", "TEST")
	if dErr != nil {
		http.Error(w, dErr.PublicError(), dErr.Status())
		log.Println(dErr.Error())
		return
	}
	char, dErr := GetCharacterFromDB(db, user.UserID)
	if dErr != nil {
		log.Println(dErr.Error())
		if err = SetNewCharacter(db, NewCharacter(user.UserID, 0, 0)); err != nil {
			http.Error(w, global.INTERNAL_SERVER_ERROR_RESPONSE, http.StatusInternalServerError)
			log.Println(err)
			return
		}
	}
	StartPlayerLoop(NewSession(user, char, ws, time.Now()))
}
