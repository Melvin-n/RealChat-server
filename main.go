package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	log.Println("Welcome to RealChat")
	client, err := fireBaseConnect()
	if err != nil {
		log.Fatalf("Unable to connect to firestore: %s", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/signup", signUp).Methods("POST")
	log.Println("Server listening on port 8080...")
	http.ListenAndServe(":8080", r)
	defer client.Close()
}
