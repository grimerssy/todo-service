package repository

import (
	"context"
	"database/sql"

	"github.com/grimerssy/todo-service/internal/core"
)

type Authorization interface {
	CreateUser(ctx context.Context, user core.User) (uint, error)
	GetUserId(ctx context.Context, username string, password string) (uint, error)
}

type Todo interface {
	Create(ctx context.Context, userId uint, todo core.Todo) (uint, error)
	GetById(ctx context.Context, userId uint, todoId uint) (core.Todo, error)
	GetByCompletion(ctx context.Context, userId uint, completed bool) ([]core.Todo, error)
	GetAll(ctx context.Context, userId uint) ([]core.Todo, error)
	Update(ctx context.Context, userId uint, todo core.Todo) error
	Patch(ctx context.Context, userId uint, todo core.Todo) error
	DeleteById(ctx context.Context, userId uint, todoId uint) error
	DeleteByCompletion(ctx context.Context, userId uint, completed bool) error
}

type Repository struct {
	Authorization
	Todo
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Todo:          NewTodoPostgres(db),
	}
}
