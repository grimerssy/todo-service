package auth

import (
	"context"
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

func (s *JWT) GenerateToken(ctx context.Context, userID any) (string, error) {
	res := make(chan func() (string, error), 1)

	go func() {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claimsJWT{
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(s.tokenTTL).Unix(),
				IssuedAt:  time.Now().Unix(),
			},
			userID,
		})

		tokenStr, err := token.SignedString([]byte(s.signingString))
		if err != nil {
			res <- func() (string, error) {
				return "", fmt.Errorf("could not sign jwt token: %s", err.Error())
			}
			return
		}

		res <- func() (string, error) {
			return tokenStr, nil
		}
		return
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case fn := <-res:
		return fn()
	}
}

func (s *JWT) ParseToken(ctx context.Context, accessToken string) (any, error) {
	res := make(chan func() (any, error), 1)

	go func() {
		token, err := jwt.ParseWithClaims(accessToken, &claimsJWT{}, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}

			return []byte(s.signingString), nil
		})
		if err != nil {
			res <- func() (any, error) {
				return nil, fmt.Errorf("could not parse jwt token: %s", err.Error())
			}
			return
		}

		claims, ok := token.Claims.(*claimsJWT)

		if !ok {
			err := errors.New(fmt.Sprintf("token claims are not of type %T", claims))
			res <- func() (any, error) {
				return nil, fmt.Errorf("could not cast token claims: %s", err.Error())
			}
			return
		}

		res <- func() (any, error) {
			return claims.UserID, nil
		}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case fn := <-res:
		return fn()
	}
}
