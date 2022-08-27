package main

import (
	"log"

	"cloud.google.com/go/firestore"
)

//global variables
var db *firestore.Client

func main() {
	log.Println("Welcome to RealChat")
	var err error
	db, err = fireBaseConnect()
	if err != nil {
		log.Fatalf("Unable to connect to firestore: %s", err)
	}
	router()
	//defer db.Close()
}
