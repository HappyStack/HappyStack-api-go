package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

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

func main() {
	initKeys()
	router := NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}

func initKeys() {
	SignKey, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatal("Error reading private key")
		return
	}
	VerifyKey, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		log.Fatal("Error reading public key")
		return
	}

	fmt.Println(SignKey)
	fmt.Println(VerifyKey)
}
