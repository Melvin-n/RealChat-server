package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/melvin-n/realchat/models"
)

func signUp(w http.ResponseWriter, r *http.Request) {
	var newUser models.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Panic("Unable to parse signUp request")
	}

	fmt.Println(newUser)
}
