package main

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
)

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
	//handlers
	http.HandleFunc("/fetchNumber", handleGetFetchNumber)
	log.Println("Server started...")
	err = http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("Starting Server failed: " + err.Error())
	}
}

func handleGetFetchNumber(w http.ResponseWriter, r *http.Request) {
	log.Println("started fetch number request...")
	query := r.URL.Query()
	number := query["number"][0]
	if number == "" {
		http.Error(w, "empty input", http.StatusBadRequest)
		log.Println("error: empty input!")
		return
	}

	log.Println("incoming number is: " + number)

	w.Write([]byte(number))
	log.Println("answered fetch number request successfully")
}
