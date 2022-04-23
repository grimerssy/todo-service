package core

import (
	"time"
)

type User struct {
	ID        uint
	FirstName string
	LastName  string
	Email     string
	Username  string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserCredentials struct {
	ID       uint
	Username string
	Password string
}

type UserRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}
