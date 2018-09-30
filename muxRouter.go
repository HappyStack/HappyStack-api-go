package main

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type MuxRouter struct {
	mux *mux.Router
}

func NewMuxRouter() *MuxRouter {
	return &MuxRouter{mux: mux.NewRouter().StrictSlash(true)}
}

func (router *MuxRouter) registerRoutes(routes []Route) {
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)
		router.mux.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
}

func (router *MuxRouter) start() error {
	return http.ListenAndServe(":8080", router.mux)
}

// URL Helper
func (router *MuxRouter) userIDForRequest(r Request) (int, error) {
	req, _ := r.(HttpRequest)
	userIDString := mux.Vars(req.httpr)["userId"]
	return strconv.Atoi(userIDString)
}

// URL Helper
func (router *MuxRouter) itemIDForRequest(r Request) (int, error) {
	req, _ := r.(HttpRequest)
	itemIDString := mux.Vars(req.httpr)["itemId"]
	return strconv.Atoi(itemIDString)
}
