package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"encoding/json"
	"io"
	"time"
	jwt "github.com/dgrijalva/jwt-go"
)


var VerifyKey, SignKey []byte

type Token struct {
	Token string `json:"token"`
}

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

const (
	privateKeyPath = "../keys/app.rsa"
	publicKeyPath  = "../keys/app.rsa.pub"
)

/* App Code */
type App struct {
	database Database
	router Router
	encryptionService EncryptionService
 }

func (app *App)run() {

	initKeys() 
	fmt.Println(SignKey)
	fmt.Println(VerifyKey)

	app.router.registerRoutes(app.routes())
	log.Fatal(app.router.start())

	defer app.database.close()
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

//List
func (app *App)list(w http.ResponseWriter, r *http.Request) {
	userIDToShow, _ := app.router.userIDForRequest(r)

	// Tell the client to expect json
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Explicitely set status code
	w.WriteHeader(http.StatusOK)
	dbItems := app.database.itemsFor(userIDToShow)

	if err := json.NewEncoder(w).Encode(dbItems); err != nil {
		panic(err)
	}
}

// Create
func (app *App)itemsCreate(w http.ResponseWriter, r *http.Request) {

	userIDToShow, _ := app.router.userIDForRequest(r)

	var item item

	// Parse the body and use LimitReader to prevent from attacks (big requests).
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	// Try to parse the JSON body into an item.
	if err := json.Unmarshal(body, &item); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	item.userId = userIDToShow
	newItem, err := app.database.create(item)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(newItem); err != nil {
		panic(err)
	}
}

// Update
func (app *App)itemsUpdate(w http.ResponseWriter, r *http.Request) {

	// token, err := jwt.ParseFromRequest(req, func(token *jwt.Token) 
	// log.Printf("Error signing the token %v\n", token)

	// TODO: Check this belongs to the currently connected user.
	itemIDToShow, _ := app.router.itemIDForRequest(r)

	// Parse the body and use LimitReader to prevent from attacks (big requests).
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	var item item
	// Try to parse the JSON body into an item.
	if err := json.Unmarshal(body, &item); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	item.Id = itemIDToShow

	newItem, err := app.database.update(item)
	if err != nil {
		log.Printf("Error signing the token %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(newItem); err != nil {
		panic(err)
	}
}

//Show
func (app *App)show(w http.ResponseWriter, r *http.Request) {
	itemIDToShow, _ := app.router.itemIDForRequest(r)
	itemToShow := app.database.read(itemIDToShow)
	json.NewEncoder(w).Encode(itemToShow)
}

// Delete
func (app *App)delete(w http.ResponseWriter, r *http.Request) {
	// Todo make sure this is the logged in user that deletes
	// his own item.
	//userID, _ := userIDForRequest(r)
	itemIDToDelete, _ := app.router.itemIDForRequest(r)

	if app.database.delete(itemIDToDelete) != nil {
		json.NewEncoder(w).Encode("DOES NOT EXIST")
	}
}

// Signup

func (app *App)signupHandler(w http.ResponseWriter, r *http.Request) {

}

// Login

func (app *App)loginHandler(w http.ResponseWriter, r *http.Request) {
	var user UserCredentials
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "Error in request")
		return
	}

	fmt.Println(user.Username, user.Password)

	password, err := app.database.passwordForUserEmail(user.Username)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "Error in request")
		return
	}

	// Here validate those are valid credentials.
	
	if !app.encryptionService.comparePasswords(password, user.Password) {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "Wrong credentials")
		return
	}

	// If so then generate auth token.
	signer := jwt.New(jwt.SigningMethodHS256)

	claims := make(jwt.MapClaims)
	claims["usename"] = user.Username
	claims["exp"] = time.Now().Add(time.Minute * 20).Unix()
	signer.Claims = claims
	tokenString, err := signer.SignedString(SignKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error while signing the token")
		log.Printf("Error signing the token %v\n", err)
	}

	token := Token{Token: tokenString}
	json.NewEncoder(w).Encode(token)
}
