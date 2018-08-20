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
		"/items",
		list,
	},
	Route{
		"ItemsCreate",
		"POST",
		"/items",
		itemsCreate,
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
		"/items/{itemId}",
		delete,
	},
}
