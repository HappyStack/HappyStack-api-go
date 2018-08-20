package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

//List
func list(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(fakeItems())

	// Tell the client to expect json
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Explicitely set status code
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(fakeItems()); err != nil {
		panic(err)
	}
}

//Show
func show(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	for _, item := range fakeItems() {
		if strconv.Itoa(item.ID) == vars["itemId"] {
			jsonItem, _ := json.Marshal(item)
			w.Write(jsonItem)
			break
		}
	}
}

func fakeItems() []item {
	var items []item
	item1 := item{
		Name:        "Vitamin D",
		Dosage:      "2000 UI",
		TakenToday:  true,
		ServingSize: 2,
		ServingType: pill}
	item2 := item{Name: "Magnesium"}
	item3 := item{Name: "Zinc"}
	items = append(items, item1, item2, item3)
	return items
}
