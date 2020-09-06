package main

import (
	global "./global"
	"database/sql"
	"github.com/gorilla/websocket"
	_ "github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
	//Error:             nil,
}

var counter = -1

var users = [][]string{{"TEST", "TEST"}, {"TEST2", "TEST2"}}

func StartWebSocketConnection(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	//Check if db connection is available
	counter++
	log.Println("counter: " + strconv.FormatInt(int64(counter), 10))
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
	user, dErr := GetUserFromDB(db, users[counter][0], users[counter][1])
	if dErr != nil {
		http.Error(w, dErr.PublicError(), dErr.Status())
		log.Println(dErr.Error())
		return
	}

	bomber := NewBomberman(user.UserID, 0, 0)

	StartPlayerLoop(NewSession(user, bomber, ws, time.Now()))
	counter--
}
