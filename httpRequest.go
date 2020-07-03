package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

type HttpRequest struct {
	httpr *http.Request
}

func (r HttpRequest) item() (item, error) {
	// Parse the body and use LimitReader to prevent from attacks (big requests).
	body, err := ioutil.ReadAll(io.LimitReader(r.httpr.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.httpr.Body.Close(); err != nil {
		panic(err)
	}

	// Try to parse the JSON body into an item.
	var item item
	err = json.Unmarshal(body, &item)
	return item, err
}

func (r HttpRequest) user() (User, error) {
	var user User
	err := json.NewDecoder(r.httpr.Body).Decode(&user)
	return user, err
}
