package main

import (
	"encoding/json"
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
func (r HttpResponse) send(stuff interface{}) {
	r.setContentType()
	if err := json.NewEncoder(r.httpr).Encode(stuff); err != nil {
		panic(err)
	}
}

func (r HttpResponse) sendEmpty() {
	r.httpr.WriteHeader(http.StatusNoContent)
}

type HttpRequest struct {
	httpr *http.Request
}
