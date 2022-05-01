package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type ConfigJWT struct {
	TokenMinutes  time.Duration
	SigningString string
}

type JWT struct {
	tokenTTL      time.Duration
	signingString string
}

type claimsJWT struct {
	jwt.StandardClaims
	UserID any
}

func NewJWT(cfg ConfigJWT) *JWT {
	return &JWT{
		tokenTTL:      cfg.TokenMinutes * time.Minute,
		signingString: cfg.SigningString,
	}
}

func (s *JWT) GenerateToken(userID any) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claimsJWT{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(s.tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		userID,
	})

	accessToken, err := token.SignedString([]byte(s.signingString))
	if err != nil {
		return "", fmt.Errorf("could not sign jwt token: %s", err.Error())
	}

	return accessToken, nil
}

func (s *JWT) ParseToken(accessToken string) (any, error) {
	keyFunc := func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(s.signingString), nil
	}

	token, err := jwt.ParseWithClaims(accessToken, &claimsJWT{}, keyFunc)
	if err != nil {
		return nil, fmt.Errorf("could not parse jwt token: %s", err.Error())
	}

	claims, ok := token.Claims.(*claimsJWT)
	if !ok {
		err := errors.New(fmt.Sprintf("token claims are not of type %T", claims))
		return nil, fmt.Errorf("could not cast token claims: %s", err.Error())
	}

	return claims.UserID, nil
}
