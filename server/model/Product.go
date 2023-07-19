package model

import (
	"time"
)

type Product struct {
	Id          int       `json:"Id"`
	CategoryId  int       `json:"CategoryId"`
	Category    string    `json:"Category"`
	Title       string    `json:"Title"`
	Price       int       `json:"Price"`
	Thumbnail   string    `json:"Thumbnail"`
	Description string    `json:"Description"`
	CreatedAt   time.Time `json:"CreatedAt"`
	UpdatedAt   time.Time `json:"UpdatedAt"`
}

type ConditionSelected struct {
	Min      int
	Max      int
	Title    string
	Category int
}
