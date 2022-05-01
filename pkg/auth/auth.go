package auth

type Authenticator interface {
	GenerateToken(userID any) (string, error)
	ParseToken(accessToken string) (any, error)
}
