package service

import (
	"context"

	"github.com/grimerssy/todo-service/internal/core"
)

type Services struct {
	AuthService
	UserService
	TodoService
}

type AuthService interface {
	GenerateToken(ctx context.Context, userReq core.UserRequest) (string, error)
	ParseToken(ctx context.Context, tokenStr string) (interface{}, error)
}

type UserService interface {
	Create(ctx context.Context, userReq core.UserRequest) error
	GetUserId(ctx context.Context, userReq core.UserRequest) (interface{}, error)
}

type TodoService interface {
	Create(ctx context.Context, userID interface{}, todoReq core.TodoRequest) error
	GetByID(ctx context.Context, userID interface{}, todoID interface{}) (core.TodoResponse, error)
	GetByCompletion(ctx context.Context, userID interface{}, completed bool) ([]core.TodoResponse, error)
	GetAll(ctx context.Context, userID interface{}) ([]core.TodoResponse, error)
	UpdateByID(ctx context.Context, userID interface{}, todoID interface{}, todoReq core.TodoRequest) error
	PatchByID(ctx context.Context, userID interface{}, todoID interface{}, todoReq core.TodoRequest) error
	DeleteByID(ctx context.Context, userID interface{}, todoID interface{}) error
	DeleteByCompletion(ctx context.Context, userID interface{}, completed bool) error
}

type Hasher interface {
	Hash(ctx context.Context, password string) (string, error)
	CompareHashAndPassword(ctx context.Context, hash string, password string) bool
}

type Encoder interface {
	Encode(ctx context.Context, id uint) (interface{}, error)
	Decode(ctx context.Context, encoded interface{}) (uint, error)
}
