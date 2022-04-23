package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncoderHashids_Encode(t *testing.T) {
	const (
		id     = 1
		salt   = "salt"
		length = 2
	)
	cfg := ConfigHashids{
		Salts: map[string]string{
			"todo": salt,
		},
		HashLength: length,
	}

	enc, err := NewEncoderHashids(cfg, TodoKey)
	require.NoError(t, err)

	hash, err := enc.Encode(context.Background(), id)
	require.NoError(t, err)

	tests := []struct {
		name     string
		input    uint
		want     interface{}
		compare  assert.ComparisonAssertionFunc
		errAsser assert.ErrorAssertionFunc
	}{
		{
			name:     "ok",
			input:    id,
			want:     hash,
			compare:  assert.Equal,
			errAsser: assert.NoError,
		},
		{
			name:     "wrong id",
			input:    0,
			want:     hash,
			compare:  assert.NotEqual,
			errAsser: assert.NoError,
		},
	}
	for _, tt := range tests {
		got, err := enc.Encode(context.Background(), tt.input)
		tt.compare(t, tt.want, got)
		tt.errAsser(t, err)
	}
}

func TestEncoderHashids_Decode(t *testing.T) {
	const (
		id     = 1
		salt   = "salt"
		length = 2
	)
	cfg := ConfigHashids{
		Salts: map[string]string{
			"todo": salt,
		},
		HashLength: length,
	}

	enc, err := NewEncoderHashids(cfg, TodoKey)
	require.NoError(t, err)

	hash, err := enc.Encode(context.Background(), id)
	require.NoError(t, err)

	tests := []struct {
		name      string
		input     interface{}
		want      uint
		compare   assert.ComparisonAssertionFunc
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name:      "ok",
			input:     hash,
			want:      id,
			compare:   assert.Equal,
			errAssert: assert.NoError,
		},
		{
			name:      "wrong id",
			input:     hash,
			want:      0,
			compare:   assert.NotEqual,
			errAssert: assert.NoError,
		},
		{
			name:      "invalid hash",
			input:     "",
			want:      0,
			compare:   assert.Equal,
			errAssert: assert.Error,
		},
	}
	for _, tt := range tests {
		got, err := enc.Decode(context.Background(), tt.input)
		tt.compare(t, tt.want, got)
		tt.errAssert(t, err)
	}
}
