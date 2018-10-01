package main

type AuthService interface {
	init()
	tokenFor(user User) (Token, error)
}

type Token struct {
	Token string `json:"token"`
}
