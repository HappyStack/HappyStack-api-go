package main

// Request abstracts a request.
type Request interface {
	item() (item, error)
	userCredentials() (UserCredentials, error)
}
