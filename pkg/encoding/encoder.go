package encoding

import (
	"context"
)

type Encoder interface {
	EncodeID(ctx context.Context, id uint) (interface{}, error)
	DecodeID(ctx context.Context, encoded interface{}) (uint, error)
}
