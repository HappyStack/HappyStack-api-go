package main

import "time"

type servingType string

const (
	scoop servingType = "scoop"
	pill  servingType = "pill"
	drop  servingType = "drop"
)

type item struct {
	Id          int         `json:"id", db:"item_id"`
	userId      int         `json:"userId", db:"user_id"`
	Name        string      `json:"name", db:"name"`
	Dosage      string      `json:"dosage"`
	TakenToday  bool        `json:"takenToday"`
	ServingSize int         `json:"servingSize"`
	ServingType servingType `json:"servingType"` // scoop, pill, drop Todo use enum later
	Timing      time.Time   `json:"timing"`
}

type Items []item
