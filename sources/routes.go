package main

import (
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func (app *App)routes() Routes {
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
	Route{
		"ItemsUpdate",
		"PUT",
		"/users/{userId}/items/{itemId}",
		app.itemsUpdate,
	},
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
	Route{
		"Signup",
		"POST",
		"/signup",
		app.signupHandler,
	},
	Route{
		"Login",
		"POST",
		"/login",
		app.loginHandler,
	},
}
}
