package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/melvin-n/realchat/models"
	"golang.org/x/crypto/bcrypt"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE")
}

func router() {
	r := mux.NewRouter()
	//TODO: create cors middleware
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

		return
	}

	if len(newUser.Username) < 4 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 5)
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
		Username: newUser.Username,
		Email:    newUser.Email,  //TODO: check email validity
		Password: hashedPassword, //TODO: check password validity
	})
	if err != nil {
		log.Printf("Unable to add to DB: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	enableCors(&w)
	response := make(map[string]string)
	response["message"] = fmt.Sprintf("Successfully created user %s", newUser.Username)
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshalling response")
		return
	}
	w.Write(jsonResponse)
	return
}

func login(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
	}
	var userRequestDetails models.User
	err := json.NewDecoder(r.Body).Decode(&userRequestDetails)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	fmt.Println(userRequestDetails.Username)
	ctx := context.Background()
	iter := db.Collection("users").Where("Username", "==", userRequestDetails.Username).Documents(ctx)

	matches, err := iter.GetAll()

	if len(matches) > 1 {
		log.Println("Error: Multiple accounts with same username")
		return
	}
	var userDBDetails models.User
	err = matches[0].DataTo(&userDBDetails)
	if err != nil {
		log.Println("Unable to retrieve user password")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(userDBDetails.Password), []byte(userRequestDetails.Password))
	if err != nil {
		log.Printf("Password does not match %s:\n", err.Error())
		return
	}

	authCookie := http.Cookie{
		Name:    "user",
		Value:   userRequestDetails.Username,
		Expires: time.Now().AddDate(0, 0, 1),
		MaxAge:  86400,
		Path:    "/",
	}

	http.SetCookie(w, &authCookie)
	w.WriteHeader(http.StatusOK)
	err = authCookie.Valid()
	if err != nil {
		log.Printf("Cookie not valid: %s", err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	response := make(map[string]string)
	response["message"] = fmt.Sprintf("Successfully logged in as %s", userRequestDetails.Username)
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshalling response")
		return
	}

	w.Write(jsonResponse)
	log.Printf("Successfully logged in as %s", userRequestDetails.Username)
	return
}

func checkForDuplicates(ctx context.Context, collection, field, value string) {
	iter := db.Collection(collection).Where(field, "==", value)
	matches, err := iter.Documents(ctx).GetAll()
	if err != nil {
		log.Fatal(err.Error())
	}

	if len(matches) > 0 {
		log.Fatalf("%s is already attached to an account", field)
		return
	}
}
