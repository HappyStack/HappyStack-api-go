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

var routes = Routes{
	Route{
		"ItemsIndex",
		"GET",
		"/users/{userId}/items",
		list,
	},
	Route{
		"ItemsCreate",
		"POST",
		"/users/{userId}/items",
		itemsCreate,
	},
	Route{
		"ItemsUpdate",
		"PUT",
		"/users/{userId}/items/{itemId}",
		itemsUpdate,
	},
	Route{
		"ItemsShow",
		"GET",
		"/items/{itemId}",
		show,
	},
	Route{
		"ItemsDelete",
		"DELETE",
		"/users/{userId}/items/{itemId}",
		delete,
	},
	Route{
		"Login",
		"POST",
		"/login",
		loginHandler,
	},
}
