package auth

import (
	"context"
)

type Authenticator interface {
	GenerateToken(ctx context.Context, userID any) (string, error)
	ParseToken(ctx context.Context, accessToken string) (any, error)
}
