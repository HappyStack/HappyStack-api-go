package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		jsonItems, _ := json.Marshal(fakeItems())
		w.Write(jsonItems)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
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

type servingType string

const (
	scoop servingType = "scoop"
	pill  servingType = "pill"
	drop  servingType = "drop"
)

type item struct {
	Name        string      `json:"name"`
	Dosage      string      `json:"dosage"`
	TakenToday  bool        `json:"takenToday"`
	ServingSize int         `json:"servingSize"`
	ServingType servingType `json:"servingType"` // scoop, pill, drop Todo use enum later
	Timing      time.Time   `json:"timing"`
}
