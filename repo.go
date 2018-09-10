package main

import "fmt"

var currentId int

var items Items

func init() {

	// Seed Data
	item1 := item{
		Name:        "Vitamin D",
		Dosage:      "2000 UI",
		TakenToday:  true,
		ServingSize: 2,
		ServingType: pill}
	item2 := item{Name: "Magnesium"}
	item3 := item{Name: "Zinc"}
	repoCreateItem(item1)
	repoCreateItem(item2)
	repoCreateItem(item3)
}

func repoFindItem(id int) item {
	for _, i := range items {
		if i.Id == id {
			return i
		}
	}
	// return empty item if not found
	return item{}
}

func repoCreateItem(i item) item {
	currentId++
	i.Id = currentId
	items = append(items, i)
	return i
}

func repoDestroyItem(id int) error {
	for i, item := range items {
		if item.Id == id {
			items = append(items[:i], items[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Could not find Todo with id of %d to delete", id)
}
