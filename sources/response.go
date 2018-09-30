package main

// Response represents the response abstraction.
type Response interface {
	send(interface{}, ResponseStatus)
	sendError(error, ResponseStatus)
}

// ResponseStatus represents the different response types.
type ResponseStatus int

// ResponseStatus enumeration
const (
	OK                  ResponseStatus = 0
	Created             ResponseStatus = 1
	NoContent           ResponseStatus = 2
	BadRequest          ResponseStatus = 3
	Unauthorized        ResponseStatus = 4
	Forbidden           ResponseStatus = 5
	NotFound            ResponseStatus = 6
	UnprocessableEntity ResponseStatus = 7
	InternalServerError ResponseStatus = 8
	NotImplemented      ResponseStatus = 9
)
