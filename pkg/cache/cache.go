package cache

type Cache interface {
	SetValue(key, val any)
	GetValue(key any) any
	RemoveValue(key any)
}

type TodoCacheKey struct {
	UserID uint
	Args   any
}
