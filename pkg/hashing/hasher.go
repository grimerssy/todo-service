package hashing

import (
	"context"
)

type Hasher interface {
	GenerateHash(ctx context.Context, password string) (string, error)
	CompareHashAndPassword(ctx context.Context, hash string, password string) bool
}
