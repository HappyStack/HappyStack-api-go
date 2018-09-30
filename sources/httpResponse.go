package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

// HttpResponse Conforms to "Response" interface
type HttpResponse struct {
	httpr http.ResponseWriter
}

func (r HttpResponse) setContentType() {
	r.httpr.Header().Set("Content-Type", "application/json; charset=UTF-8")
}
func (r HttpResponse) setStatusOK() {
	r.httpr.WriteHeader(http.StatusOK)
}
func (r HttpResponse) setStatusCreated() {
	r.httpr.WriteHeader(http.StatusCreated)
}

func (r HttpResponse) setStatusUnprocessableEntity() {
	r.httpr.WriteHeader(http.StatusUnprocessableEntity)
}

func (r HttpResponse) send(stuff interface{}) {
	r.setContentType()
	if err := json.NewEncoder(r.httpr).Encode(stuff); err != nil {
		panic(err)
	}
}

func (r HttpResponse) sendEmpty() {
	r.httpr.WriteHeader(http.StatusNoContent)
}

func (r HttpResponse) sendError(e error) {
	r.setContentType()
	r.httpr.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(r.httpr).Encode(e); err != nil {
		panic(err)
	}
}

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
