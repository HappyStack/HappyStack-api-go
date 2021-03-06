package main

import (
	"fmt"
	"log"
)

type User struct {
	Id       int    `json:"id", db:"user_id"`
	Username string `json:"username"`
	Password string `json:"password"`
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

//List
func (app *App) list(res Response, req Request) {
	userIDToShow, _ := app.router.userIDForRequest(req)
	dbItems := app.database.itemsFor(userIDToShow)
	res.send(dbItems, OK)
}

// Create
func (app *App) itemsCreate(res Response, req Request) {

	if !app.authService.hasAuthorization(req) {
		res.send("Needs authentication", Forbidden)
		return
	}

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

	if !app.authService.isAuthorizedForUserId(userID, req) {
		res.send("invalid token", Forbidden)
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

	if !app.authService.hasAuthorization(req) {
		res.send("Needs authentication", Forbidden)
		return
	}

	// Parse item.
	item, err := req.item()
	if err != nil {
		res.sendError(err, UnprocessableEntity)
		return
	}
	itemID, _ := app.router.itemIDForRequest(req)
	item.Id = itemID

	// Fetch item from Database to see correponding user_id.
	// the one in the request could be a fake one !
	dbitem, err := app.database.read(itemID)
	if err != nil {
		res.sendError(err, NotFound)
	}

	if !app.authService.isAuthorizedForUserId(dbitem.userId, req) {
		res.send("invalid token", Forbidden)
		return
	}

	// Update it
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
	item, err := app.database.read(itemID)
	if err != nil {
		res.sendError(err, NotFound)
	}
	res.send(item, OK)
}

// Delete
func (app *App) delete(res Response, req Request) {

	if !app.authService.hasAuthorization(req) {
		res.send("Needs authentication", Forbidden)
		return
	}

	itemID, _ := app.router.itemIDForRequest(req)
	item, err := app.database.read(itemID)
	if err != nil {
		res.sendError(err, NotFound)
	}
	itemBelongsToUserID := item.userId

	if !app.authService.isAuthorizedForUserId(itemBelongsToUserID, req) {
		res.send("invalid token", Forbidden)
		return
	}

	err = app.database.delete(itemID)
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
	user, err := req.user()
	if err != nil {
		res.sendError(err, UnprocessableEntity)
		return
	}

	fmt.Println(user.Username, user.Password)

	dbUser, err := app.database.userMatchingEmail(user.Username)
	if err != nil {
		res.sendError(err, BadRequest)
		return
	}

	// Here validate those are valid credentials.
	if !app.encryptionService.comparePasswords(dbUser.Password, user.Password) {
		res.send("Wrong credentials", Forbidden)
		return
	}

	// If so then generate auth token.
	token, err := app.authService.tokenFor(dbUser)
	if err != nil {
		res.sendError(err, BadRequest)
	}
	res.send(token, OK)
}
