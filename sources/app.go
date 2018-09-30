package main

import (
	"fmt"
	"log"
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

type Request interface {
	item() (item, error)
	userCredentials() (UserCredentials, error)
}

//List
func (app *App) list(res Response, req Request) {
	userIDToShow, _ := app.router.userIDForRequest(req)
	dbItems := app.database.itemsFor(userIDToShow)
	res.send(dbItems, OK)
}

// Create
func (app *App) itemsCreate(res Response, req Request) {
	item, err := req.item()
	if err != nil {
		res.sendError(err, UnprocessableEntity)
		return
	}

	userID, err := app.router.userIDForRequest(req)
	if err != nil {
		res.sendError(err, BadRequest)
		return
	}
	item.userId = userID
	newItem, err := app.database.create(item)
	if err != nil {
		res.sendError(err, BadRequest)
		return
	}
	res.send(newItem, Created)
}

// Update
func (app *App) itemsUpdate(res Response, req Request) {
	// TODO: Auth Check this belongs to the currently connected user.

	// Parse item.
	item, err := req.item()
	if err != nil {
		res.sendError(err, UnprocessableEntity)
		return
	}
	itemID, _ := app.router.itemIDForRequest(req)
	item.Id = itemID

	// Create it
	updatedItem, err := app.database.update(item)
	if err != nil {
		res.sendError(err, BadRequest)
		return
	}
	res.send(updatedItem, OK)
}

//Show
func (app *App) show(res Response, req Request) {
	itemID, _ := app.router.itemIDForRequest(req)
	item := app.database.read(itemID)
	res.send(item, OK)
}

// Delete
func (app *App) delete(res Response, req Request) {
	itemID, _ := app.router.itemIDForRequest(req)
	err := app.database.delete(itemID)
	if err != nil {
		res.send("DOES NOT EXIST", NotFound)
	} else {
		res.send("", NoContent)
	}
}

// Signup

func (app *App) signupHandler(res Response, req Request) {

}

// Login

func (app *App) loginHandler(res Response, req Request) {

	// Parse credentials.
	user, err := req.userCredentials()
	if err != nil {
		res.sendError(err, UnprocessableEntity)
		return
	}

	fmt.Println(user.Username, user.Password)

	password, err := app.database.passwordForUserEmail(user.Username)
	if err != nil {
		res.sendError(err, BadRequest)
		return
	}

	// Here validate those are valid credentials.
	if !app.encryptionService.comparePasswords(password, user.Password) {
		res.send("Wrong credentials", Forbidden)
		return
	}

	// If so then generate auth token.
	token, err := app.authService.tokenFor(user)
	if err != nil {
		res.sendError(err, BadRequest)
		// w.WriteHeader(http.StatusInternalServerError)
		// fmt.Fprintln(w, "Error while signing the token")
		// log.Printf("Error signing the token %v\n", err)
	}
	res.send(token, OK)
}
