package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

// URL Helper
func userIDForRequest(r *http.Request) (int, error) {
	userIDString := mux.Vars(r)["userId"]
	return strconv.Atoi(userIDString)
}

//List
func list(w http.ResponseWriter, r *http.Request) {
	userIDToShow, _ := userIDForRequest(r)

	// Tell the client to expect json
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Explicitely set status code
	w.WriteHeader(http.StatusOK)
	dbItems := happyStackDatabase.allItemsForUserId(userIDToShow)

	if err := json.NewEncoder(w).Encode(dbItems); err != nil {
		panic(err)
	}
}

// Create
func itemsCreate(w http.ResponseWriter, r *http.Request) {

	userIDToShow, _ := userIDForRequest(r)

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
	newItem, err := happyStackDatabase.createItem(item)
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
func itemsUpdate(w http.ResponseWriter, r *http.Request) {

	// TODO: Check this belongs to the currently connected user.
	itemIDToShow, _ := itemIDForRequest(r)

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

	newItem, err := happyStackDatabase.updateItem(item)
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
func show(w http.ResponseWriter, r *http.Request) {
	itemIDToShow, _ := itemIDForRequest(r)
	itemToShow := happyStackDatabase.findItem(itemIDToShow)
	json.NewEncoder(w).Encode(itemToShow)
}

// Delete
func delete(w http.ResponseWriter, r *http.Request) {
	// Todo make sure this is the logged in user that deletes
	// his own item.
	//userID, _ := userIDForRequest(r)
	itemIDToDelete, _ := itemIDForRequest(r)

	if happyStackDatabase.destroyItem(itemIDToDelete) != nil {
		json.NewEncoder(w).Encode("DOES NOT EXIST")
	}
}

// URL Helper
func itemIDForRequest(r *http.Request) (int, error) {
	itemIDString := mux.Vars(r)["itemId"]
	return strconv.Atoi(itemIDString)
}

// Signup

func signupHandler(w http.ResponseWriter, r *http.Request) {

}

// Login

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var user UserCredentials
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "Error in request")
		return
	}

	fmt.Println(user.Username, user.Password)

	password, err := happyStackDatabase.passwordForUserEmail(user.Username)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "Error in request")
		return
	}

	// Here validate those are valid credentials.
	wrongCredentials := (user.Password != password)
	if wrongCredentials {
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
