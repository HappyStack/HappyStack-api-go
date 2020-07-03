package main

type AuthService interface {
	init()
	tokenFor(user User) (Token, error)
	hasAuthorization(req Request) bool
	isAuthorizedForUserId(userId int, req Request) bool
}

type Token struct {
	Token string `json:"token"`
}
