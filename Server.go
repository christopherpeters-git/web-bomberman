package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

//database
const (
	DB_USERNAME = "root"
	DB_PASSWORD = "soe2020"
	DB_URL      = "localhost:3306"
	DB_NAME     = "krass"
)

//handler url's
const (
	POST_SAVEPICTURE = "/uploadImage"
	WEBSOCKET_TEST   = "/ws-test/"
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

	db, err = sql.Open("mysql", DB_NAME+":"+DB_PASSWORD+"@tcp("+DB_URL+")/"+DB_NAME)
	if err != nil {
		log.Fatal("Database connection failed: " + err.Error())
	}
	defer db.Close()

	//handlers
	http.HandleFunc(POST_SAVEPICTURE, handleUploadImage)
	http.HandleFunc(WEBSOCKET_TEST, handleWebsocketEndpoint)
	log.Println("Server started...")
	err = http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("Starting Server failed: " + err.Error())
	}
}

func handleWebsocketEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Println("websocket started...")
	if err := StartWebSocketConnection(w, r); err != nil {
		log.Println(err.Error())
		return
	}
	log.Println("websocket connected...")
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
