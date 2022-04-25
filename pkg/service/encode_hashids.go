package service

import (
	"context"
	"errors"

	"github.com/speps/go-hashids"
)

const (
	UserKey = "user"
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

	return &EncoderHashids{
		hashID: hID,
	}, err
}

func (e *EncoderHashids) Encode(ctx context.Context, id uint) (interface{}, error) {
	res := make(chan func() (interface{}, error), 1)

	go func() {
		hash, err := e.hashID.EncodeInt64([]int64{int64(id)})
		res <- func() (interface{}, error) {
			return hash, err
		}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case f := <-res:
		return f()
	}
}

func (e *EncoderHashids) Decode(ctx context.Context, encoded interface{}) (uint, error) {
	res := make(chan func() (uint, error), 1)

	go func() {
		hash, ok := encoded.(string)
		if !ok {
			res <- func() (uint, error) {
				return 0, errors.New("given value could be converted to string")
			}
			return
		}

		ids, err := e.hashID.DecodeInt64WithError(hash)
		if err != nil {
			res <- func() (uint, error) {
				return 0, err
			}
			return
		}

		if len(ids) == 0 {
			res <- func() (uint, error) {
				return 0, errors.New("invalid hash")
			}
			return
		}
		id := uint(ids[0])

		res <- func() (uint, error) {
			return id, nil
		}
	}()

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case f := <-res:
		return f()
	}
}
