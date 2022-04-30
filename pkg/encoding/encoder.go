package encoding

import (
	"context"
)

type Encoder interface {
	EncodeID(ctx context.Context, id uint) (any, error)
	DecodeID(ctx context.Context, encoded any) (uint, error)
}
