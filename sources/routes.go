package main

import (
	"net/http"
)

type Route struct {
	Name    string
	Method  string
	Pattern string
	// HandlerFunc http.HandlerFunc
	HandlerFunc MyHandlerFunc
}

type MyHandlerFunc func(Response, Request)

// ServeHTTP calls f(w, r).
func (f MyHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wrappedResponse := HttpResponse{httpr: w}
	wrappedRequest := HttpRequest{httpr: r}
	f(wrappedResponse, wrappedRequest)
}

type Routes []Route

func (app *App) routes() Routes {
	return Routes{
		Route{
			"ItemsIndex",
			"GET",
			"/users/{userId}/items",
			app.list,
		},
		Route{
			"ItemsCreate",
			"POST",
			"/users/{userId}/items",
			app.itemsCreate,
		},
		// Route{
		// 	"ItemsUpdate",
		// 	"PUT",
		// 	"/users/{userId}/items/{itemId}",
		// 	app.itemsUpdate,
		// },
		Route{
			"ItemsShow",
			"GET",
			"/users/{userId}/items/{itemId}",
			app.show,
		},
		Route{
			"ItemsDelete",
			"DELETE",
			"/users/{userId}/items/{itemId}",
			app.delete,
		},
		// Route{
		// 	"Signup",
		// 	"POST",
		// 	"/signup",
		// 	app.signupHandler,
		// },
		// Route{
		// 	"Login",
		// 	"POST",
		// 	"/login",
		// 	app.loginHandler,
		// },
	}
}
