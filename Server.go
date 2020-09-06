package main

import (
	glo "./global"
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

//database
const (
	DB_USERNAME = "test"
	DB_PASSWORD = "test123"
	DB_URL      = "connectolai.de:3306"
	DB_NAME     = "testdb"
)

//handler url's
const (
	POST_SAVEPICTURE             = "/uploadImage"
	WEBSOCKET_TEST               = "/ws-test/"
	GET_FETCH_ACTIVE_CONNECTIONS = "/fetchConnections/"
	POST_LOGIN                   = "/login"
	POST_REGISTER                = "/register"
)

var db *sql.DB

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

	//go UpdateClients()

	//handlers
	http.HandleFunc(POST_LOGIN, handleLogin)
	http.HandleFunc(POST_SAVEPICTURE, handleUploadImage)
	http.HandleFunc(WEBSOCKET_TEST, handleWebsocketEndpoint)
	http.HandleFunc(GET_FETCH_ACTIVE_CONNECTIONS, handleFetchActiveConnections)
	log.Println("Server started...")
	err = http.ListenAndServe(":2100", nil)
	if err != nil {
		log.Fatal("Starting Server failed: " + err.Error())
	}
}

func handleFetchActiveConnections(w http.ResponseWriter, r *http.Request) {
	log.Println("handling fetch active connections request started...")
	w.Write([]byte(AllConnectionsAsString()))
	log.Println("handling fetch active connections request ended...")
}

func handleWebsocketEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Println("handling websocket started...")
	StartWebSocketConnection(w, r, db)
	log.Println("handling websocket ended...")
}

func handleUploadImage(w http.ResponseWriter, r *http.Request) {
	log.Println("Upload started...")
	//Parsing ??? Maxsize = 10mb
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Println("Parsing failed: " + err.Error())
	}
	//Retrieving
	file, handler, err := r.FormFile("imageFile")
	if err != nil {
		log.Println("Retrieving failed: " + err.Error())
		return
	}
	defer file.Close()

	log.Println("Uploaded File: ", handler.Filename)
	log.Println("File size: ", handler.Size)
	log.Println("MIME Header: ", handler.Header)

	//Writing
	//TO-DO: Change TempFile func
	tempFile, err := ioutil.TempFile("temp-images", "upload-*.png")
	if err != nil {
		log.Println("Writing failed: " + err.Error())
		return
	}

	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err.Error())
	}

	tempFile.Write(fileBytes)

	log.Println(w, "Successfully Uploaded!")
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("Receiving Loginrequest...")
	err := db.Ping()
	if err != nil {
		log.Println("Database connection failed" + err.Error())
		return
	}
	err = r.ParseForm()
	if err != nil {
		http.Error(w, glo.INTERNAL_SERVER_ERROR_RESPONSE, http.StatusInternalServerError)
		log.Println("Parsing Form failed for some reason" + err.Error())
		return
	}

	username := r.FormValue("usernameInput")
	password := r.FormValue("passwordInput")

	if !isStringLegal(username) {
		log.Println("Parsed username contains illegal chars or is empty!")
		return
	}
	if !isStringLegal(password) {
		log.Println("Parsed password contains illegal chars or is empty!")
		return
	}
	user, httpErr := GetUserFromDB(db, username, password)
	if user == nil || httpErr != nil {
		//todo print error correctly
		log.Println("Couldnt get user from database for some reason")
		return
	}
	userAsJson, err := json.MarshalIndent(user, "", "    ")
	if err != nil {
		log.Println("Marshaling failed" + err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(userAsJson)
	log.Println("Login successfully handled")

}

func handleRegister(w http.ResponseWriter, r *http.Request) {

}
