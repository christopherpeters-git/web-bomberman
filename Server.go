package main

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/gorilla/websocket"
	"io/ioutil"
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
	http.HandleFunc("/upload", handleUploadImage)

	log.Println("Server started...")
	err = http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("Starting Server failed: " + err.Error())
	}
}

func handleUploadImage(w http.ResponseWriter, r *http.Request) {
	log.Println("Upload started...")
	//Parsing ??? Maxsize = 10mb
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Println("Parsing failed: " + err.Error())
	}
	//Retrieve
	file, handler, err := r.FormFile("imageFile")
	if err != nil {
		log.Println("Retrieving failed: " + err.Error())
		return
	}
	defer file.Close()

	log.Println("Uploaded File: ", handler.Filename)
	log.Println("File size: ", handler.Size)
	//What is MIME Header?
	log.Println("MIME Header: ", handler.Header)

	//write to temp
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
