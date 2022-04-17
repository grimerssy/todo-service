package core

import (
	"time"
)

type Todo struct {
	Id          uint
	Title       string
	Description string
	Completed   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type TodoCreate struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

type TodoUpdate struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type TodoResponse struct {
	Id          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Completed   bool   `json:"completed"`
}
