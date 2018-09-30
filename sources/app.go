package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthService interface {
	init()
	tokenFor(user UserCredentials) (Token, error)
}

/* App Code */
type App struct {
	database          Database
	router            Router
	encryptionService EncryptionService
	authService       AuthService
}

func (app *App) run() {
	app.authService.init()
	app.router.registerRoutes(app.routes())
	log.Fatal(app.router.start())
	defer app.database.close()
}

type Response interface {
	setContentType()
	setStatusOK()
	send(stuff interface{})
	sendEmpty()
}

type Request interface {
	// userID() (int, error)
}

//List
func (app *App) list(res Response, req Request) {
	userIDToShow, _ := app.router.userIDForRequest(req)
	res.setContentType()
	res.setStatusOK()
	dbItems := app.database.itemsFor(userIDToShow)
	res.send(dbItems)
}

// Create
func (app *App) itemsCreate(w http.ResponseWriter, r *http.Request) {

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
func (app *App) itemsUpdate(w http.ResponseWriter, r *http.Request) {

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
func (app *App) show(res Response, req Request) {
	itemID, _ := app.router.itemIDForRequest(req)
	item := app.database.read(itemID)
	res.send(item)
}

func (app *App) oldShow(w http.ResponseWriter, r *http.Request) {
	itemIDToShow, _ := app.router.itemIDForRequest(r)
	itemToShow := app.database.read(itemIDToShow)
	json.NewEncoder(w).Encode(itemToShow)
}

// Delete
func (app *App) delete(res Response, req Request) {
	itemID, _ := app.router.itemIDForRequest(req)
	err := app.database.delete(itemID)
	if err != nil {
		res.send("DOES NOT EXIST")
	} else {
		res.sendEmpty()
	}
}

// Signup

func (app *App) signupHandler(w http.ResponseWriter, r *http.Request) {

}

// Login

func (app *App) loginHandler(w http.ResponseWriter, r *http.Request) {
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
	token, err := app.authService.tokenFor(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error while signing the token")
		log.Printf("Error signing the token %v\n", err)
	}

	json.NewEncoder(w).Encode(token)
}
