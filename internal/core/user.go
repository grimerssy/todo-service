package core

import (
	"time"
)

type User struct {
	Id           uint
	FirstName    string
	LastName     string
	Email        string
	Username     string
	Password     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	AuthorizedAt time.Time
}

type UserSignUp struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

type UserSignIn struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
