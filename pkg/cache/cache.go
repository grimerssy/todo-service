package cache

type Cache interface {
	SetValue(key, val interface{})
	GetValue(key interface{}) interface{}
}
