package cache

type Cache interface {
	SetValue(key, val interface{})
	GetValue(key interface{}) interface{}
	RemoveValue(key interface{})
}

type TodoCacheKey struct {
	UserID uint
	Args   interface{}
}
