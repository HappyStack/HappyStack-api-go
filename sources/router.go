package main

type Router interface {
	registerRoutes(routes []Route)
	start() error
	userIDForRequest(r Request) (int, error)
	itemIDForRequest(r Request) (int, error)
}
