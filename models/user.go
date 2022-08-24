package models

type User struct {
	Username       string `json:"username"`
	Id             string `json:"id"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
}
