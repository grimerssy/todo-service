package hashing

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHashBcrypt_Hash(t *testing.T) {
	const (
		cost     = 10
		password = "password"
	)
	cfg := ConfigBcrypt{
		Cost: cost,
	}

	h := NewBcrypt(cfg)

	tests := []struct {
		name        string
		input       string
		bcryptInput string
		compare     assert.ComparisonAssertionFunc
		errAssert   assert.ErrorAssertionFunc
	}{
		{
			name:        "ok",
			input:       password,
			bcryptInput: password,
			compare:     assert.Equal,
			errAssert:   assert.NoError,
		},
		{
			name:        "different passwords",
			input:       password,
			bcryptInput: "",
			compare:     assert.NotEqual,
			errAssert:   assert.NoError,
		},
	}
	for _, tt := range tests {
		hash, err := h.GenerateHash(tt.input)
		tt.errAssert(t, err)
		bytes, err := bcrypt.GenerateFromPassword([]byte(tt.bcryptInput), cfg.Cost)
		require.NoError(t, err)

		got := bcrypt.CompareHashAndPassword([]byte(hash), []byte(tt.input))
		want := bcrypt.CompareHashAndPassword(bytes, []byte(tt.input))
		tt.compare(t, want, got)
	}
}

func TestHashBcrypt_CompareHashAndPassword(t *testing.T) {
	const (
		cost     = 10
		password = "password"
	)
	cfg := ConfigBcrypt{
		Cost: cost,
	}

	h := NewBcrypt(cfg)

	hash, err := h.GenerateHash(password)
	require.NoError(t, err)

	tests := []struct {
		name     string
		hash     string
		password string
	}{
		{
			name:     "ok",
			hash:     hash,
			password: password,
		},
	}
	for _, tt := range tests {
		got := h.CompareHashAndPassword(tt.hash, tt.password)
		want := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
		assert.Equal(t, want, got)
	}
}
