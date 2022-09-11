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

func accessControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("mw")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		next.ServeHTTP(w, r)
	})
}
func router() {
	r := mux.NewRouter()
	r.Use(accessControlMiddleware)
	r.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		return
	})
	r.HandleFunc("/signup", signUp).Methods("POST")
	r.HandleFunc("/login", login).Methods("POST")
	log.Println("Server listening on port 8080...")
	http.ListenAndServe(":8080", r)
}

func signUp(w http.ResponseWriter, r *http.Request) {

	var newUser models.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusOK)

		return
	}

	if len(newUser.Username) < 4 {
		handleResponse(&w, http.StatusOK, "Username is too short - must be 4 characters or more.")
		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 5)
	if err != nil {
		log.Println("Password hashing error")
		w.WriteHeader(http.StatusOK)
	}
	hashedPassword := string(password)

	ctx := context.Background()

	//search if user or email already exists
	err = checkForDuplicates(ctx, "users", "Username", newUser.Username)
	err = checkForDuplicates(ctx, "users", "Email", newUser.Email)
	if err != nil {
		handleResponse(&w, http.StatusOK, err.Error())
	}

	_, _, err = db.Collection("users").Add(ctx, models.User{
		Username: newUser.Username,
		Email:    newUser.Email,  //TODO: check email validity
		Password: hashedPassword, //TODO: check password validity
	})
	if err != nil {
		log.Printf("Unable to add to DB: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}

	handleResponse(&w, http.StatusCreated, fmt.Sprintf("Successfully created user %s", newUser.Username))
	return
}

func login(w http.ResponseWriter, r *http.Request) {
	var userRequestDetails models.User
	err := json.NewDecoder(r.Body).Decode(&userRequestDetails)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	fmt.Println(userRequestDetails.Username)
	ctx := context.Background()
	iter := db.Collection("users").Where("Username", "==", userRequestDetails.Username).Documents(ctx)

	if userRequestDetails.Username == "" || userRequestDetails.Password == "" {
		handleResponse(&w, http.StatusOK, "Can not process blank fields")
		log.Println("Can not process blank fields")
		return
	}

	matches, err := iter.GetAll()

	if len(matches) > 1 {
		handleResponse(&w, http.StatusOK, "Error: Multiple accounts with same username")
		log.Println("Error: Multiple accounts with same username")
		return
	}

	if len(matches) == 0 {
		handleResponse(&w, http.StatusOK, "Incorrect username or password")
		log.Println("Incorrect username or password")
		return
	}
	var userDBDetails models.User
	err = matches[0].DataTo(&userDBDetails)
	if err != nil {
		log.Println("Unable to retrieve user password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(userDBDetails.Password), []byte(userRequestDetails.Password))
	if err != nil {
		handleResponse(&w, http.StatusOK, "Incorrect username or password")
		log.Printf("Password does not match %s:\n", err.Error())

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

	handleResponse(&w, http.StatusOK, fmt.Sprintf("Successfully logged in as %s", userRequestDetails.Username))

	log.Printf("Successfully logged in as %s", userRequestDetails.Username)
	return
}

func checkForDuplicates(ctx context.Context, collection, field, value string) error {
	iter := db.Collection(collection).Where(field, "==", value)
	matches, err := iter.Documents(ctx).GetAll()
	if err != nil {
		log.Fatal(err.Error())
	}

	if len(matches) > 0 {
		return fmt.Errorf("%s is already attached to an account", field)
	}

	return nil
}

func handleResponse(w *http.ResponseWriter, status int, message string) {
	(*w).WriteHeader(status)
	(*w).Header().Set("Content-Type", "application/json")
	response := make(map[string]string)
	response["message"] = message
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshalling response")
		return
	}
	(*w).Write(jsonResponse)
}
