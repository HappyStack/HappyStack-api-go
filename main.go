package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var happyStackDatabase *HappyStackDatabase

func main() {
	initKeys()

	happyStackDatabase = NewHappyStackDatabase()
	defer happyStackDatabase.closeDatabase()

	fmt.Println(SignKey)
	fmt.Println(VerifyKey)

	router := NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}

func initKeys() {
	sKey, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatal("Error reading private key")
		return
	}
	SignKey = sKey
	vKey, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		log.Fatal("Error reading public key")
		return
	}
	VerifyKey = vKey
}
