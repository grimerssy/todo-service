package encoding

type Encoder interface {
	EncodeID(id uint) (any, error)
	DecodeID(encoded any) (uint, error)
}
