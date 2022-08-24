package main

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"github.com/melvin-n/realchat/models"
	"google.golang.org/api/option"
)

func main() {
	fmt.Println("Welcome to RealChat")
	fireBaseConnect()
}

func fireBaseConnect() {
	ctx := context.Background()
	opt := option.WithCredentialsFile("secrets/realchat27-firebase-admin-SDK.json")
	config := &firebase.Config{ProjectID: "RealChat27"}

	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		log.Fatalf("error initializing firebase connection: %s", err.Error())
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Unable to connect to client: %s", err.Error())
	}

	defer client.Close()
	//TODO: fix db permissions - unable to write
	_, _, err = client.Collection("users").Add(ctx, models.User{
		Username:       "Melvin",
		Id:             "1",
		Email:          "test@mail.com",
		HashedPassword: "####",
	})
	if err != nil {
		log.Fatalf("Unable to add to DB: %s", err.Error())
	}
}
