package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "sacha"
	password = ""
	dbname   = "sacha"
)

func main() {
	initKeys()

	// TODO add DB password
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
		host, port, user, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to database!")

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
