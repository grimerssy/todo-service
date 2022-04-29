package encoding

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

type Hashids struct {
	hashID *hashids.HashID
}

func NewHashids(cfg ConfigHashids, saltKey string) (*Hashids, error) {
	data := hashids.NewData()
	data.MinLength = int(cfg.HashLength)
	data.Salt = cfg.Salts[saltKey]
	hID, err := hashids.NewWithData(data)

	return &Hashids{
		hashID: hID,
	}, err
}

func (e *Hashids) EncodeID(ctx context.Context, id uint) (interface{}, error) {
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

func (e *Hashids) DecodeID(ctx context.Context, encoded interface{}) (uint, error) {
	res := make(chan func() (uint, error), 1)

	go func() {
		hash, ok := encoded.(string)
		if !ok {
			res <- func() (uint, error) {
				return 0, errors.New("given value could not be converted to string")
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