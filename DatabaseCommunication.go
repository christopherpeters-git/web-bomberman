package main

import (
	global "./global"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"strings"
)

type User struct {
	UserID       uint64
	Username     string
	passwordHash string
	sessionID    string
}

func NewUser(userID uint64, username string, passwordHash string, sessionID string) *User {
	return &User{UserID: userID, Username: username, passwordHash: passwordHash, sessionID: sessionID}
}

type Character struct {
	UserID      uint64
	PositionX   int
	PositionY   int
	PicturePath string
}

const LETTER_BYTES = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789?!"

func (r *Character) String() string {
	return "Character: {" + strconv.FormatUint(r.UserID, 10) + " | " + strconv.FormatInt(int64(r.PositionX), 10) + " | " + strconv.FormatInt(int64(r.PositionY), 10) + " | " + r.PicturePath + "}"
}

func NewCharacter(userID uint64, positionX int, positionY int) *Character {
	return &Character{UserID: userID, PositionX: positionX, PositionY: positionY}
}

func (r *User) String() string {
	return "User: {" + strconv.FormatUint(r.UserID, 10) + " | " + r.Username + "}"
}

func GetUserFromDB(db *sql.DB, username string, password string) (*User, *global.DetailedHttpError) {
	var user User
	if err := db.Ping(); err != nil {
		return nil, global.NewDetailedHttpError(http.StatusInternalServerError, global.INTERNAL_SERVER_ERROR_RESPONSE, err.Error())
	}
	rows, err := db.Query("select * from users where username = ?", username)
	if err != nil {
		return nil, global.NewDetailedHttpError(http.StatusInternalServerError, global.INTERNAL_SERVER_ERROR_RESPONSE, err.Error())
	}
	if rows.Next() {
		if err = rows.Scan(&user.UserID, &user.Username, &user.passwordHash, &user.sessionID); err != nil {
			return nil, global.NewDetailedHttpError(http.StatusInternalServerError, global.INTERNAL_SERVER_ERROR_RESPONSE, err.Error())
		}
	} else {
		return nil, global.NewDetailedHttpError(http.StatusNotFound, "User not found", "User not found: "+username)
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.passwordHash), []byte(password)); err == nil { //TODO WRONG
		return nil, global.NewDetailedHttpError(http.StatusBadRequest, "wrong password", "wrong password")
	}
	return &user, nil
}

func GetCharacterFromDB(db *sql.DB, userId uint64) (*Character, *global.DetailedHttpError) {
	var char Character
	if err := db.Ping(); err != nil {
		return nil, global.NewDetailedHttpError(http.StatusInternalServerError, global.INTERNAL_SERVER_ERROR_RESPONSE, err.Error())
	}
	rows, err := db.Query("select * from characters where users_id = ?", userId)
	if err != nil {
		return nil, global.NewDetailedHttpError(http.StatusInternalServerError, global.INTERNAL_SERVER_ERROR_RESPONSE, err.Error())
	}
	if rows.Next() {
		if err = rows.Scan(&char.UserID, &char.PositionX, &char.PositionY, &char.PicturePath); err != nil {
			return nil, global.NewDetailedHttpError(http.StatusInternalServerError, global.INTERNAL_SERVER_ERROR_RESPONSE, err.Error())
		}
	} else {
		return nil, global.NewDetailedHttpError(http.StatusInternalServerError, global.INTERNAL_SERVER_ERROR_RESPONSE, "no character found for "+strconv.FormatUint(userId, 10))
	}
	return &char, nil
}

func SetNewCharacter(db *sql.DB, char *Character) error {
	//TODO write function
	return nil
}

func isStringLegal(str string) bool {
	if str == "" {
		return false
	}
	for _, c := range str {
		if !strings.Contains(LETTER_BYTES, strings.ToLower(string(c))) {
			return false
		}
	}
	return true
}
