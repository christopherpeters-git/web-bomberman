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

const LETTER_BYTES = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789?!"

func NewUser(userID uint64, username string, passwordHash string, sessionID string) *User {
	return &User{UserID: userID, Username: username, passwordHash: passwordHash, sessionID: sessionID}
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
