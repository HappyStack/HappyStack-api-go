package main

const (
	privateKeyPath = "keys/app.rsa"
	publicKeyPath  = "keys/app.rsa.pub"
)

var VerifyKey, SignKey []byte

type Token struct {
	Token string `json:"token"`
}

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
