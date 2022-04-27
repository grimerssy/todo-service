package service

import (
	"context"
	"errors"

	"github.com/grimerssy/todo-service/internal/core"
)

var (
	ErrTodoNotFound = errors.New("todo does not exist")
	ErrUserNotFound = errors.New("user does not exist")
)

type Services struct {
	UserService
	TodoService
}

type UserService interface {
	SignUp(ctx context.Context, userReq core.UserRequest) error
	SignIn(ctx context.Context, userReq core.UserRequest) (string, error)
	GetID(ctx context.Context, token string) (interface{}, error)
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
