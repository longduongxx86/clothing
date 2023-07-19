package model

import "time"

type Account struct {
	Id          int
	Fullname    string
	Email       string
	PhoneNumber string
	Address     string
	Password    string
	RoleId      int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
