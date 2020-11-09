package main

import (
	"database/sql"
	"github.com/gorilla/websocket"
	_ "github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  128,
	WriteBufferSize: 128,
	CheckOrigin:     func(r *http.Request) bool { return true },
	//Error:             nil,
}

func StartWebSocketConnection(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	//Check if db connection is available
	if err := db.Ping(); err != nil {
		http.Error(w, INTERNAL_SERVER_ERROR_RESPONSE, http.StatusInternalServerError)
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
	rand.Seed(time.Now().UTC().UnixNano())
	random := rand.Intn(4)
	bomber := NewBomberman(user.UserID, 0, 0, user.Username)
	if sessionRunning {
		bomber.Kill()
		findWinner()
	}
	GameMap.Fields[pixToArr(bomber.PositionX)][pixToArr(bomber.PositionY)].Player.PushBack(bomber)

	if random == 0 {
		bomber.teleportTo(19, 0, pixToArr(bomber.PositionX), pixToArr(bomber.PositionY))
	} else if random == 1 {
		bomber.teleportTo(0, 19, pixToArr(bomber.PositionX), pixToArr(bomber.PositionY))
	} else if random == 2 {
		bomber.teleportTo(19, 19, pixToArr(bomber.PositionX), pixToArr(bomber.PositionY))
	}

	StartPlayerLoop(NewSession(user, bomber, ws, time.Now()))

}
