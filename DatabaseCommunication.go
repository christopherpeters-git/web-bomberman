package main

/*
All Methods which belong to the Database Communication.
*/

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

/*
Represents a User in the Database.
*/
type User struct {
	UserID       uint64
	Username     string
	passwordHash string
	sessionID    string
	GamesWon     uint64
}

const LETTER_BYTES = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789?!"
const COOKIE_NAME = "BMCookie"

func NewUser(userID uint64, username string, passwordHash string, sessionID string) *User {
	return &User{UserID: userID, Username: username, passwordHash: passwordHash, sessionID: sessionID}
}

func (r *User) String() string {
	return "User: {" + strconv.FormatUint(r.UserID, 10) + " | " + r.Username + "}"
}

/*
Gets a User by its Username and Password.
*/
func GetUserFromDB(db *sql.DB, username string, password string) (*User, *DetailedHttpError) {
	var user User
	if err := db.Ping(); err != nil {
		return nil, NewDetailedHttpError(http.StatusInternalServerError, INTERNAL_SERVER_ERROR_RESPONSE, err.Error())
	}
	rows, err := db.Query("select * from users where username = ?", username)
	if err != nil {
		return nil, NewDetailedHttpError(http.StatusInternalServerError, INTERNAL_SERVER_ERROR_RESPONSE, err.Error())
	}
	if rows.Next() {
		if err = rows.Scan(&user.UserID, &user.Username, &user.passwordHash, &user.sessionID, &user.GamesWon); err != nil {
			return nil, NewDetailedHttpError(http.StatusInternalServerError, INTERNAL_SERVER_ERROR_RESPONSE, err.Error())
		}
	} else {
		return nil, NewDetailedHttpError(http.StatusNotFound, "User not found", "User not found: "+username)
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.passwordHash), []byte(password)); err != nil {
		return nil, NewDetailedHttpError(http.StatusBadRequest, "wrong password", "wrong password")
	}
	return &user, nil
}

/*
Gets a User by its SessionID.
*/
func GetUserFromSessionCookie(db *sql.DB, sessionId string) (*User, *DetailedHttpError) {
	var user User
	if err := db.Ping(); err != nil {
		return nil, NewDetailedHttpError(http.StatusInternalServerError, INTERNAL_SERVER_ERROR_RESPONSE, err.Error())
	}
	rows, err := db.Query("select * from users where session_id = ?", sessionId)
	if err != nil {
		return nil, NewDetailedHttpError(http.StatusInternalServerError, INTERNAL_SERVER_ERROR_RESPONSE, err.Error())
	}
	if rows.Next() {
		if err = rows.Scan(&user.UserID, &user.Username, &user.passwordHash, &user.sessionID, &user.GamesWon); err != nil {
			return nil, NewDetailedHttpError(http.StatusInternalServerError, INTERNAL_SERVER_ERROR_RESPONSE, err.Error())
		}
		return &user, nil

	}
	return nil, NewDetailedHttpError(http.StatusNotFound, "No user found for this Session-ID", "No user found for this Session-ID")
}

/*
Checks if a User with the Username-String passed already exists.
*/
func UsernameExists(db *sql.DB, username string) *DetailedHttpError {
	err := db.Ping()
	if err != nil {
		log.Println("Database connection failed" + err.Error())
		return NewDetailedHttpError(http.StatusInternalServerError, INTERNAL_SERVER_ERROR_RESPONSE, err.Error())
	}
	rows, err := db.Query("select * from users where username = ?", username)
	if err != nil {
		log.Println("Something went wrong on sql.Query" + err.Error())
		return NewDetailedHttpError(http.StatusInternalServerError, INTERNAL_SERVER_ERROR_RESPONSE, err.Error())
	}
	if rows.Next() {
		return NewDetailedHttpError(http.StatusInternalServerError, "Username is already taken", err.Error())
	}
	return nil
}

func IsStringLegal(str string) bool {
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

/*
Generates a SessionID and places a Cookie with this ID.
Placed Cookie lasts 2 days.
*/
func PlaceCookie(w http.ResponseWriter, db *sql.DB, username string) error {
	rows, err := db.Query("select * from users where username = ?", username)
	if err != nil {
		return err
	}
	sessionId, err := generateUniqueSessionId(db)
	if err != nil {
		return err
	}
	if rows.Next() {
		_, err = db.Exec("UPDATE  users set Session_Id = ? where username = ?", sessionId, username)
		if err != nil {
			return err
		}
	}
	expire := time.Now().AddDate(0, 0, 2)
	cookie := http.Cookie{
		Name:       COOKIE_NAME,
		Value:      sessionId,
		Path:       "/",
		Domain:     "localhost",
		Expires:    expire,
		RawExpires: expire.Format(time.UnixDate),
		MaxAge:     172800,
		Secure:     false,
		HttpOnly:   true,
		SameSite:   http.SameSiteLaxMode,
	}
	log.Println("placed cookie successfully")
	http.SetCookie(w, &cookie)
	return nil
}

/*
Checks if a Cookie for the User is present.
Returns nil if everything is alright.
*/
func CheckCookie(r *http.Request, db *sql.DB, user *User) *DetailedHttpError {
	cookie, err := r.Cookie(COOKIE_NAME)
	if err != nil {
		return NewDetailedHttpError(http.StatusNotFound, "No cookie found", err.Error())
	}
	rows, err := db.Query("select * from users where session_id = ?", cookie.Value)
	if err != nil {
		return NewDetailedHttpError(http.StatusNotFound, "Session_Id doesnt exists", err.Error())
	}
	if rows.Next() {
		err = rows.Scan(&user.UserID, &user.Username, &user.passwordHash, &user.sessionID, &user.GamesWon)
		if err != nil {
			return NewDetailedHttpError(http.StatusInternalServerError, INTERNAL_SERVER_ERROR_RESPONSE, err.Error())
		}
	}
	return nil
}

/*
Generates SessionID's until a unique one is found.
*/
func generateUniqueSessionId(db *sql.DB) (string, error) {
	sessionId := generateSessionId(255)
	rows, err := db.Query("select session_id from users where session_id = ?", sessionId)
	if err != nil {
		return "", err
	}
	if rows.Next() {
		return generateUniqueSessionId(db)
	}
	return sessionId, nil
}

/*
Generates a SessionID.
*/
func generateSessionId(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = LETTER_BYTES[rand.Intn(len(LETTER_BYTES))]
	}
	return string(b)
}

/*
Gets a User by its ID.
*/
func getUserByID(db *sql.DB, userID uint64) (*User, *DetailedHttpError) {
	var user User
	if err := db.Ping(); err != nil {
		return nil, NewDetailedHttpError(http.StatusInternalServerError, INTERNAL_SERVER_ERROR_RESPONSE, err.Error())
	}
	rows, err := db.Query("select * from users where Id = ?", userID)
	if err != nil {
		return nil, NewDetailedHttpError(http.StatusInternalServerError, INTERNAL_SERVER_ERROR_RESPONSE, err.Error())
	}
	if rows.Next() {
		if err = rows.Scan(&user.UserID, &user.Username, &user.passwordHash, &user.sessionID, &user.GamesWon); err != nil {
			return nil, NewDetailedHttpError(http.StatusInternalServerError, INTERNAL_SERVER_ERROR_RESPONSE, err.Error())
		}
		return &user, nil

	}
	return nil, NewDetailedHttpError(http.StatusNotFound, "No user found for this Session-ID", "No user found for this Session-ID")
}

/*
Updates the Player-Statistics in the Database.
*/
func updatePlayerStats(db *sql.DB, user User) *DetailedHttpError {
	log.Println("Updating Player Statistic..")
	err := db.Ping()
	if err != nil {
		log.Println("Database connection failed" + err.Error())
		return NewDetailedHttpError(http.StatusInternalServerError, INTERNAL_SERVER_ERROR_RESPONSE, err.Error())
	}

	rows, err := db.Query("select * from users where Id = ?", user.UserID)
	if err != nil {
		return NewDetailedHttpError(http.StatusInternalServerError, INTERNAL_SERVER_ERROR_RESPONSE, err.Error())
	}

	log.Println(user.GamesWon)

	if rows.Next() {
		_, err = db.Exec("UPDATE  users set Games_won = ? where ID = ?", user.GamesWon, user.UserID)
		if err != nil {
			return NewDetailedHttpError(http.StatusInternalServerError, INTERNAL_SERVER_ERROR_RESPONSE, err.Error())
		}
	}

	return nil
}
