package core

import (
	"time"
)

type Todo struct {
	ID          uint
	Title       string
	Description string
	Completed   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type TodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

type TodoResponse struct {
	ID          interface{} `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description,omitempty"`
	Completed   bool        `json:"completed"`
}
