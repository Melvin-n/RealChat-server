package main

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func fireBaseConnect() (*firestore.Client, error) {
	ctx := context.Background()
	opt := option.WithCredentialsFile("secrets/realchat27-firebase-adminsdk-rksx0-5720746270.json")
	config := &firebase.Config{ProjectID: "realchat27"}

	log.Println("Connecting to firebase...")
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		return nil, err
	}

	log.Println("Connecting to firestore...")
	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	//defer client.Close()
	//TODO: fix db permissions - unable to write
	if err != nil {
		log.Fatalf("Unable to add to DB: %s", err.Error())
	}

	return client, nil
}
