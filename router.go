package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/melvin-n/realchat/models"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/iterator"
)

func router() {
	r := mux.NewRouter()
	r.HandleFunc("/signup", signUp).Methods("POST")
	r.HandleFunc("/login", login).Methods("POST")
	log.Println("Server listening on port 8080...")
	http.ListenAndServe(":8080", r)
}

func signUp(w http.ResponseWriter, r *http.Request) {
	var newUser models.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Panic("Unable to parse signUp request")
	}

	if len(newUser.Username) < 4 {
		w.WriteHeader(http.StatusBadRequest)
	}

	password, err := bcrypt.GenerateFromPassword([]byte(newUser.HashedPassword), 5)
	if err != nil {
		log.Println("Password hashing error")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	hashedPassword := string(password)

	ctx := context.Background()

	//search if user or email already exists
	checkForDuplicates(ctx, "users", "Username", newUser.Username)
	checkForDuplicates(ctx, "users", "Email", newUser.Email)

	_, _, err = db.Collection("users").Add(ctx, models.User{
		Username:       newUser.Username,
		Email:          newUser.Email,  //TODO: check email validity
		HashedPassword: hashedPassword, //TODO: check password validity
	})
	if err != nil {
		log.Fatalf("Unable to add to DB: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK) //TODO: write custom status messages
	return
}

func login(w http.ResponseWriter, r *http.Request) {
	var userDetails models.User
	err := json.NewDecoder(r.Body).Decode(&userDetails)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Panic("Unable to parse signUp request")
	}
	fmt.Println(userDetails.Username)
	ctx := context.Background()
	iter := db.Collection("users").Where("Username", "==", userDetails.Username).Documents(ctx)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Println(doc.Data())
	}

}

func checkForDuplicates(ctx context.Context, collection, field, value string) {
	iter := db.Collection(collection).Where(field, "==", value)
	matches, err := iter.Documents(ctx).GetAll()
	if err != nil {
		log.Fatal(err.Error())
	}

	if len(matches) > 0 {
		log.Fatalf("%s is already attached to an account", field)
	}
}
