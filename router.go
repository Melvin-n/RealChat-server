package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/melvin-n/realchat/models"
	"golang.org/x/crypto/bcrypt"
)

type appHandler struct {
	Handler func(http.ResponseWriter, *http.Request) (int, error)
}

func (ah appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status, err := ah.Handler(w, r)
	if err != nil {
		log.Printf("HTTP: %d: %v", status, err)
	}

	switch status {
	case http.StatusNotFound:
		http.NotFound(w, r)
	case http.StatusBadRequest:
		http.Error(w, err.Error(), http.StatusBadRequest)
	case http.StatusInternalServerError:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	default:
		http.Error(w, err.Error(), status)
	}
}

func router() {
	r := mux.NewRouter()
	r.HandleFunc("/signup", signUp).Methods("POST")
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

	context := context.Background()
	_, _, err = db.Collection("users").Add(context, models.User{
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
