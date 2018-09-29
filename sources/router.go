package main

import "net/http"

type Router interface {
	registerRoutes(routes []Route)
	start() error
	userIDForRequest(r *http.Request) (int, error)
	itemIDForRequest(r *http.Request) (int, error) 
}