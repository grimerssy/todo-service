package auth

import (
	"context"
)

type Authenticator interface {
	GenerateToken(ctx context.Context, userID interface{}) (string, error)
	ParseToken(ctx context.Context, accessToken string) (interface{}, error)
}
