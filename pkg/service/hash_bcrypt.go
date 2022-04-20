package service

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type ConfigBcrypt struct {
	Cost int
}

type HashBcrypt struct {
	cost int
}

func NewHashBcrypt(cfg ConfigBcrypt) *HashBcrypt {
	return &HashBcrypt{cost: cfg.Cost}
}

func (h *HashBcrypt) Hash(ctx context.Context, password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)

	return string(hash), err
}

func (h *HashBcrypt) CompareHashAndPassword(ctx context.Context, hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}
