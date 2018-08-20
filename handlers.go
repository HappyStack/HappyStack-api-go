package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

//List
func list(w http.ResponseWriter, r *http.Request) {

	// Tell the client to expect json
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Explicitely set status code
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(items); err != nil {
		panic(err)
	}
}

// Create
func itemsCreate(w http.ResponseWriter, r *http.Request) {

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

	t := repoCreateItem(item)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		panic(err)
	}
}

//Show
func show(w http.ResponseWriter, r *http.Request) {
	itemIDToShow, _ := itemIDForRequest(r)
	itemToShow := repoFindItem(itemIDToShow)
	json.NewEncoder(w).Encode(itemToShow)
}

// Delete
func delete(w http.ResponseWriter, r *http.Request) {
	itemIDToDelete, _ := itemIDForRequest(r)

	if repoDestroyItem(itemIDToDelete) != nil {
		json.NewEncoder(w).Encode("DOES NOT EXIST")
	}
}

// URL Helper
func itemIDForRequest(r *http.Request) (int, error) {
	itemIDString := mux.Vars(r)["itemId"]
	return strconv.Atoi(itemIDString)
}
