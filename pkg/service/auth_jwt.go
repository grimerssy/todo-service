package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/grimerssy/todo-service/internal/core"
)

type ConfigJWT struct {
	minutesTTL uint
	secret     string
}

type AuthJWT struct {
	tokenTTL    time.Duration
	secretJWT   string
	userService UserService
}

type claimsJWT struct {
	jwt.StandardClaims
	UserId interface{} `json:"userId"`
}

func NewAuthJWT(cfg ConfigJWT, user UserService) *AuthJWT {
	return &AuthJWT{
		tokenTTL:    time.Duration(cfg.minutesTTL) * time.Minute,
		secretJWT:   cfg.secret,
		userService: user,
	}
}

func (s *AuthJWT) GenerateToken(ctx context.Context, userSI core.UserSignIn) (string, error) {
	userID, err := s.userService.GetUserId(ctx, userSI)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claimsJWT{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(s.tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		userID,
	})

	return token.SignedString([]byte(s.secretJWT))
}

func (s *AuthJWT) ParseToken(ctx context.Context, tokenStr string) (interface{}, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &claimsJWT{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(s.secretJWT), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*claimsJWT)

	if !ok {
		return nil, errors.New(fmt.Sprintf("token claims are not of type %T", claims))
	}

	return claims.UserId, nil
}
