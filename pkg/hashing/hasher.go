package hashing

type Hasher interface {
	GenerateHash(password string) (string, error)
	CompareHashAndPassword(hash string, password string) bool
}
