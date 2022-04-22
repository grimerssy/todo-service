package repository

import (
	"context"

	"github.com/grimerssy/todo-service/internal/core"
)

type User interface {
	Create(ctx context.Context, user core.User) error
	GetCredentialsByUsername(ctx context.Context, username string) (core.UserAuth, error)
}

type Todo interface {
	Create(ctx context.Context, userId uint, todo core.Todo) (uint, error)
	GetByID(ctx context.Context, userId uint, todoId uint) (core.Todo, error)
	GetByCompletion(ctx context.Context, userId uint, completed bool) ([]core.Todo, error)
	GetAll(ctx context.Context, userId uint) ([]core.Todo, error)
	Update(ctx context.Context, userId uint, todoId uint, todo core.Todo) error
	Patch(ctx context.Context, userId uint, todoId uint, todo core.Todo) error
	DeleteByID(ctx context.Context, userId uint, todoId uint) error
	DeleteByCompletion(ctx context.Context, userId uint, completed bool) error
}

type Repository struct {
	User
	Todo
}

func NewRepository(user User, todo Todo) *Repository {
	return &Repository{
		User: user,
		Todo: todo,
	}
}
