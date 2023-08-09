package main

import (
	"log"
	"math/rand"
	"time"
)

type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}
type Account struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Number    int64     `json:"number"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewAccount(firstName, lastName string) *Account {
	log.Println("inside new account")
	return &Account{
		// ID:        rand.Intn(10000),
		FirstName: firstName,
		LastName:  lastName,
		Number:    int64(rand.Intn(100000)),
		CreatedAt: time.Now().UTC(),
	}
}
