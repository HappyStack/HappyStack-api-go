package main

type AuthService interface {
	init()
	tokenFor(user UserCredentials) (Token, error)
}

type Token struct {
	Token string `json:"token"`
}
