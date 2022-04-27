package hashing

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type ConfigBcrypt struct {
	Cost int
}

type Bcrypt struct {
	cost int
}

func NewBcrypt(cfg ConfigBcrypt) *Bcrypt {
	return &Bcrypt{
		cost: cfg.Cost,
	}
}

func (h *Bcrypt) GenerateHash(ctx context.Context, password string) (string, error) {
	res := make(chan func() (string, error), 1)

	go func() {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)

		res <- func() (string, error) {
			return string(hash), err
		}
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case f := <-res:
		return f()
	}
}

func (h *Bcrypt) CompareHashAndPassword(ctx context.Context, hash string, password string) bool {
	res := make(chan bool, 1)

	go func() {
		err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

		res <- err == nil
	}()

	select {
	case <-ctx.Done():
		return false
	case match := <-res:
		return match
	}
}
