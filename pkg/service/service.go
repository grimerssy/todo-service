package service

import (
	"context"

	"github.com/grimerssy/todo-service/internal/core"
)

type Service struct {
	AuthenticationService
}

type AuthenticationService interface {
	GenerateToken(ctx context.Context, userReq core.UserRequest) (string, error)
	ParseToken(ctx context.Context, tokenStr string) (interface{}, error)
}

type UserService interface {
	Create(ctx context.Context, userReq core.UserRequest) error
	GetUserId(ctx context.Context, userReq core.UserRequest) (interface{}, error)
}

type Hasher interface {
	Hash(ctx context.Context, password string) (string, error)
	CompareHashAndPassword(ctx context.Context, hash string, password string) bool
}

type Encoder interface {
	Encode(ctx context.Context, id uint) (interface{}, error)
	Decode(ctx context.Context, encoded interface{}) (uint, error)
}
