package hashing

import (
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

func (h *Bcrypt) GenerateHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	return string(hash), err
}

func (h *Bcrypt) CompareHashAndPassword(hash string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
