package repository

import (
	"context"

	"github.com/grimerssy/todo-service/internal/core"
)

type Repositories struct {
	UserRepository
	TodoRepository
}

type UserRepository interface {
	Create(ctx context.Context, user core.User) error
	GetCredentialsByUsername(ctx context.Context, username string) (core.UserCredentials, error)
}

type TodoRepository interface {
	Create(ctx context.Context, userID uint, todo core.Todo) error
	GetByID(ctx context.Context, userID uint, todoID uint) (core.Todo, error)
	GetByCompletion(ctx context.Context, userID uint, completed bool) ([]core.Todo, error)
	GetAll(ctx context.Context, userID uint) ([]core.Todo, error)
	UpdateByID(ctx context.Context, userID uint, todoID uint, todo core.Todo) error
	PatchByID(ctx context.Context, userID uint, todoID uint, todo core.Todo) error
	DeleteByID(ctx context.Context, userID uint, todoID uint) error
	DeleteByCompletion(ctx context.Context, userID uint, completed bool) error
}
