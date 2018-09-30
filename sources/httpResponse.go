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

func (r HttpResponse) setStatus(s ResponseStatus) {
	httpStatus := httpStatusForResponseStatus(s)
	r.httpr.WriteHeader(httpStatus)
}

func httpStatusForResponseStatus(s ResponseStatus) int {
	if s == OK {
		return http.StatusOK
	}
	if s == Created {
		return http.StatusCreated
	}
	if s == NoContent {
		return http.StatusNoContent
	}
	if s == BadRequest {
		return http.StatusBadRequest
	}
	if s == Unauthorized {
		return http.StatusUnauthorized
	}
	if s == Forbidden {
		return http.StatusForbidden
	}
	if s == NotFound {
		return http.StatusNotFound
	}
	if s == UnprocessableEntity {
		return http.StatusUnprocessableEntity
	}
	if s == InternalServerError {
		return http.StatusInternalServerError
	}
	if s == NotImplemented {
		return http.StatusNotImplemented
	}

	return http.StatusOK
}

func (r HttpResponse) setStatusForbidden() {
	r.httpr.WriteHeader(http.StatusForbidden)
}

func (r HttpResponse) setStatusBadRequest() {
	r.httpr.WriteHeader(http.StatusBadRequest)
}

func (r HttpResponse) send(stuff interface{}, s ResponseStatus) {
	r.setStatus(s)
	if s != NoContent {
		r.setContentType()
		if err := json.NewEncoder(r.httpr).Encode(stuff); err != nil {
			panic(err)
		}
	}
}

func (r HttpResponse) sendError(e error, s ResponseStatus) {
	r.setContentType()
	r.setStatus(s)
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

func (r HttpRequest) userCredentials() (UserCredentials, error) {
	var user UserCredentials
	err := json.NewDecoder(r.httpr.Body).Decode(&user)
	return user, err
}
