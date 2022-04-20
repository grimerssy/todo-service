package repository

import (
	"context"

	"github.com/grimerssy/todo-service/internal/core"
)

type Authentication interface {
	CreateUser(ctx context.Context, user core.User) (uint, error)
	GetUserId(ctx context.Context, username string, password string) (uint, error)
}

type Todo interface {
	Create(ctx context.Context, userId uint, todo core.Todo) (uint, error)
	GetById(ctx context.Context, userId uint, todoId uint) (core.Todo, error)
	GetByCompletion(ctx context.Context, userId uint, completed bool) ([]core.Todo, error)
	GetAll(ctx context.Context, userId uint) ([]core.Todo, error)
	Update(ctx context.Context, userId uint, todoId uint, todo core.Todo) error
	Patch(ctx context.Context, userId uint, todoId uint, todo core.Todo) error
	DeleteById(ctx context.Context, userId uint, todoId uint) error
	DeleteByCompletion(ctx context.Context, userId uint, completed bool) error
}

type Repository struct {
	Authentication
	Todo
}

func NewRepository(auth Authentication, todo Todo) *Repository {
	return &Repository{
		Authentication: auth,
		Todo:           todo,
	}
}
