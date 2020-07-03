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

func (r HttpResponse) setStatus(s ResponseStatus) {
	httpStatus := httpStatusForResponseStatus(s)
	r.httpr.WriteHeader(httpStatus)
}

func httpStatusForResponseStatus(s ResponseStatus) int {
	switch s {
	case OK:
		return http.StatusOK
	case Created:
		return http.StatusCreated
	case NoContent:
		return http.StatusNoContent
	case BadRequest:
		return http.StatusBadRequest
	case Unauthorized:
		return http.StatusUnauthorized
	case Forbidden:
		return http.StatusForbidden
	case NotFound:
		return http.StatusNotFound
	case UnprocessableEntity:
		return http.StatusUnprocessableEntity
	case InternalServerError:
		return http.StatusInternalServerError
	case NotImplemented:
		return http.StatusNotImplemented
	default:
		return http.StatusOK
	}
}

func (r HttpResponse) send(stuff interface{}, s ResponseStatus) {
	r.setContentType()
	r.setStatus(s)
	if s != NoContent {
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
