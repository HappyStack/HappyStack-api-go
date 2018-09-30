package main

import (
	"encoding/json"
	"fmt"
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
	setStatusCreated()
	setStatusUnprocessableEntity()
	send(stuff interface{})
	sendError(error)
	sendEmpty()
}

type Request interface {
	item() (item, error)
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
func (app *App) itemsCreate(res Response, req Request) {
	item, err := req.item()
	if err != nil {
		res.setStatusUnprocessableEntity()
		res.sendError(err)
		return
	}

	userID, err := app.router.userIDForRequest(req)
	if err != nil {
		res.sendError(err)
		return
	}
	item.userId = userID
	newItem, err := app.database.create(item)
	if err != nil {
		res.sendError(err)
		return
	}
	res.setStatusCreated()
	res.send(newItem)
}

// Update
func (app *App) itemsUpdate(res Response, req Request) {

	// Parse item.
	item, err := req.item()
	if err != nil {
		res.setStatusUnprocessableEntity()
		res.sendError(err)
		return
	}

	// TODO: Check this belongs to the currently connected user.
	itemID, _ := app.router.itemIDForRequest(req)
	item.Id = itemID

	newItem, err := app.database.update(item)
	if err != nil {
		//TO SET STATUS code?
		res.sendError(err)
		return
	}
	res.setStatusOK()
	res.send(newItem)
}

//Show
func (app *App) show(res Response, req Request) {
	itemID, _ := app.router.itemIDForRequest(req)
	item := app.database.read(itemID)
	res.send(item)
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

func (app *App) signupHandler(res Response, req Request) {

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
