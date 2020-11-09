package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

//database

//handler url's
const (
	WEBSOCKET_TEST               = "/ws-test/"
	GET_FETCH_ACTIVE_CONNECTIONS = "/fetchConnections/"
	POST_LOGIN                   = "/login"
	POST_REGISTER                = "/register"
	GET_FETCH_USER_ID            = "/fetchUserId"
	GET_SET_READY                = "/setReady"
	REQUEST_TIMEOUT_MILLIS       = 500
)

var db *sql.DB
var ipTimers map[uint64]bool = make(map[uint64]bool)

func main() {
	//Creates a log file
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	f, err := os.OpenFile("Server.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	http.Handle("/", http.FileServer(http.Dir("frontend/")))

	db, err = sql.Open("mysql", DB_USERNAME+":"+DB_PASSWORD+"@tcp("+DB_URL+")/"+DB_NAME)
	if err != nil {
		log.Fatal("Database connection failed: " + err.Error())
	}
	defer db.Close()

	//Initialize Game
	initGame()

	//handlers
	http.HandleFunc(POST_REGISTER, handleRegister)
	http.HandleFunc(POST_LOGIN, handleLogin)
	http.HandleFunc(WEBSOCKET_TEST, handleWebsocketEndpoint)
	http.HandleFunc(GET_FETCH_ACTIVE_CONNECTIONS, handleFetchActiveConnections)
	http.HandleFunc(GET_FETCH_USER_ID, handleGetUserID)
	http.HandleFunc(GET_SET_READY, handleGetSetReady)
	log.Println("Server started...")
	err = http.ListenAndServe(":2100", nil)
	if err != nil {
		log.Fatal("Starting Server failed: " + err.Error())
	}
}

/*
converts ip to uint64, returns on err a 0 and an error from ParseUint
*/
func ipToInt(ip string) (uint64, error) {
	log.Println(ip)
	ipString := strings.Split(strings.ReplaceAll(ip, ".", ""), ":")[0]
	log.Println(ipString)
	ipInt, err := strconv.ParseUint(ipString, 10, 64)
	if err != nil {
		return 0, err
	}
	return ipInt, nil
}

/*
Checks if ip is allowed to do another request and starts a timer if allowed
*/
func checkIpTimer(ip string) bool {
	if strings.HasPrefix(ip, "[::1]") { //For a local connection
		ip = "127.0.0.1"
	}
	ipInt, err := ipToInt(ip)
	if err != nil {
		log.Println(err)
		return false
	}
	if !ipTimers[ipInt] {
		//Starts the timer for the ip
		go func() {
			ipTimers[ipInt] = true
			time.Sleep(time.Millisecond * REQUEST_TIMEOUT_MILLIS)
			ipTimers[ipInt] = false
		}()
		return true
	}
	return false
}

func handleGetSetReady(w http.ResponseWriter, r *http.Request) {
	log.Println("handling handleGetSetReady request started...")
	if !checkIpTimer(r.RemoteAddr) {
		log.Println("not allowed")
		w.WriteHeader(http.StatusTooManyRequests)
		return
	}
	var user User
	if dErr := CheckCookie(r, db, &user); dErr != nil {
		log.Println(dErr.Error())
		http.Error(w, dErr.PublicError(), dErr.Status())
		return
	}
	Connections[user.UserID].Bomber.PlayerReady = !Connections[user.UserID].Bomber.PlayerReady
	msg := "nrdy"
	if Connections[user.UserID].Bomber.PlayerReady {
		msg = "rdy"
		StartGameIfPlayersReady()
	}
	w.Write([]byte(msg))
	log.Println("handling handleGetSetReady request ended...")
}

func handleGetUserID(w http.ResponseWriter, r *http.Request) {
	log.Println("handling handleGetUserID request started...")
	var user User
	if dErr := CheckCookie(r, db, &user); dErr != nil {
		log.Println(dErr.Error())
		http.Error(w, dErr.PublicError(), dErr.Status())
		return
	}
	w.Write([]byte(strconv.FormatUint(user.UserID, 10)))
	log.Println("handling handleGetUserID request ended...")
}

func handleFetchActiveConnections(w http.ResponseWriter, r *http.Request) {
	log.Println("handling fetch active Connections request started...")
	w.Write([]byte(AllConnectionsAsString()))
	log.Println("handling fetch active Connections request ended...")
}

func handleWebsocketEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Println("handling websocket started...")
	StartWebSocketConnection(w, r, db)
	log.Println("handling websocket ended...")
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("Receiving Loginrequest...")
	err := r.ParseForm()
	if err != nil {
		http.Error(w, INTERNAL_SERVER_ERROR_RESPONSE, http.StatusInternalServerError)
		log.Println("Parsing Form failed for some reason" + err.Error())
		return
	}

	username := r.FormValue("usernameInput")
	password := r.FormValue("passwordInput")

	if !IsStringLegal(username) {
		log.Println("Parsed username contains illegal chars or is empty")
		return
	}
	if !IsStringLegal(password) {
		log.Println("Parsed password contains illegal chars or is empty")
		return
	}
	user, httpErr := GetUserFromDB(db, username, password)
	if httpErr != nil {
		http.Error(w, httpErr.PublicError(), httpErr.Status())
		log.Println(httpErr.Error())
		return
	}
	userAsJson, err := json.MarshalIndent(user, "", "    ")
	if err != nil {
		log.Println("Marshaling failed" + err.Error())
		return
	}
	err = PlaceCookie(w, db, username)
	if err != nil {
		http.Error(w, INTERNAL_SERVER_ERROR_RESPONSE, http.StatusInternalServerError)
		log.Println(err.Error())
	}
	w.WriteHeader(http.StatusOK)
	w.Write(userAsJson)
	log.Println("Login successfully handled")

}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	log.Println("Receiving Registerrequest...")
	err := r.ParseForm()
	if err != nil {
		http.Error(w, INTERNAL_SERVER_ERROR_RESPONSE, http.StatusInternalServerError)
		log.Println("Parsing Form failed for some reason" + err.Error())
		return
	}
	username := r.FormValue("usernameInput")
	password := r.FormValue("passwordInput")

	if !IsStringLegal(username) {
		log.Println("Parsed username contains illegal chars or is empty")
		http.Error(w, INTERNAL_SERVER_ERROR_RESPONSE, http.StatusInternalServerError)
		return
	}
	if !IsStringLegal(password) {
		log.Println("Parsed password contains illegal chars or is empty")
		http.Error(w, INTERNAL_SERVER_ERROR_RESPONSE, http.StatusInternalServerError)
		return
	}
	httpErr := UsernameExists(db, username)
	if httpErr != nil {
		http.Error(w, httpErr.PublicError(), httpErr.Status())
		log.Println(httpErr.Error())
		return
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Println("encrypting password failed" + err.Error())
		http.Error(w, INTERNAL_SERVER_ERROR_RESPONSE, http.StatusInternalServerError)
		return
	}
	log.Printf("User created: username: %s passwordhash: %s", username, string(passwordHash))
	//Create user in database
	_, err = db.Exec("INSERT INTO users (Username,PasswordHash)\nValues(?,?)", username, string(passwordHash))
	if err != nil {
		log.Println("Creating entry in database failed" + err.Error())
		http.Error(w, INTERNAL_SERVER_ERROR_RESPONSE, http.StatusInternalServerError)
		return
	}
	err = PlaceCookie(w, db, username)
	if err != nil {
		http.Error(w, INTERNAL_SERVER_ERROR_RESPONSE, http.StatusInternalServerError)
		log.Println(err.Error())
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Neuer Account angelegt"))
	log.Println("handle register sucessfull")

}

func handleCookie(w http.ResponseWriter, r *http.Request) {
	var user User
	httpErr := CheckCookie(r, db, &user)
	if httpErr != nil {
		http.Error(w, httpErr.PublicError(), httpErr.Status())
		log.Println(httpErr.Error())
		return
	}
	userAsJson, err := json.MarshalIndent(user, "", "    ")
	if err != nil {
		log.Println("Marshaling failed" + err.Error())
		http.Error(w, INTERNAL_SERVER_ERROR_RESPONSE, http.StatusInternalServerError)
		return
	}
	w.Write(userAsJson)
	log.Println("cookie handled successfully")
}
