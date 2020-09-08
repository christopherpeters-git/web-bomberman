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
	cookie, err := r.Cookie(COOKIE_NAME)
	if err != nil {
		http.Error(w, "No Cookie found", http.StatusNotFound)
		log.Print(err)
		return
	}

	user, dErr := GetUserFromSessionCookie(db, cookie.Value)
	if dErr != nil {
		http.Error(w, dErr.PublicError(), dErr.Status())
		log.Println(dErr.Error())
		return
	}

	bomber := NewBomberman(user.UserID, 0, 0, user.Username)

	StartPlayerLoop(NewSession(user, bomber, ws, time.Now()))

}
