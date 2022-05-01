package encoding

import (
	"errors"

	"github.com/speps/go-hashids"
)

type cfgKey string

const (
	UserKey cfgKey = "user"
	TodoKey cfgKey = "todo"
)

type ConfigHashids struct {
	Salts       map[cfgKey]string
	HashLengths map[cfgKey]uint
}

type Hashids struct {
	hashID *hashids.HashID
}

func NewHashids(cfg ConfigHashids, cfgKey cfgKey) (*Hashids, error) {
	data := hashids.NewData()
	data.MinLength = int(cfg.HashLengths[cfgKey])
	data.Salt = cfg.Salts[cfgKey]
	hID, err := hashids.NewWithData(data)

	return &Hashids{
		hashID: hID,
	}, err
}

func (e *Hashids) EncodeID(id uint) (any, error) {
	return e.hashID.EncodeInt64([]int64{int64(id)})
}

func (e *Hashids) DecodeID(encoded any) (uint, error) {
	hash, ok := encoded.(string)
	if !ok {
		return 0, errors.New("given value could not be converted to string")
	}

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
