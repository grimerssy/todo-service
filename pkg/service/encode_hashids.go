package service

import (
	"context"
	"errors"

	"github.com/speps/go-hashids"
)

const (
	TodoKey = "todo"
)

type ConfigHashids struct {
	Salts      map[string]string
	HashLength uint
}

type EncoderHashids struct {
	hashID *hashids.HashID
}

func NewEncoderHashids(cfg ConfigHashids, saltKey string) (*EncoderHashids, error) {
	data := hashids.NewData()
	data.MinLength = int(cfg.HashLength)
	data.Salt = cfg.Salts[saltKey]
	hID, err := hashids.NewWithData(data)

	return &EncoderHashids{hashID: hID}, err
}

func (e *EncoderHashids) Encode(ctx context.Context, id uint) (string, error) {
	hash, err := e.hashID.EncodeInt64([]int64{int64(id)})
	return hash, err
}

func (e *EncoderHashids) Decode(ctx context.Context, hash string) (uint, error) {
	ids, err := e.hashID.DecodeInt64WithError(hash)
	if err != nil {
		return 0, err
	}

	if len(ids) == 0 {
		return 0, errors.New("invalid hash")
	}
	id := uint(ids[0])

	return id, nil
}
