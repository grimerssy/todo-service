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
	MinutesTTL uint
	Secret     string
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
		tokenTTL:    time.Duration(cfg.MinutesTTL) * time.Minute,
		secretJWT:   cfg.Secret,
		userService: user,
	}
}

func (s *AuthJWT) GenerateToken(ctx context.Context, userReq core.UserRequest) (string, error) {
	userID, err := s.userService.GetUserId(ctx, userReq)
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

	tokenStr, err := token.SignedString([]byte(s.secretJWT))
	if err != nil {
		return "", fmt.Errorf("could not sign jwt token: %s", err.Error())
	}

	return tokenStr, nil
}

func (s *AuthJWT) ParseToken(ctx context.Context, tokenStr string) (interface{}, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &claimsJWT{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(s.secretJWT), nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not parse jwt token: %s", err.Error())
	}

	claims, ok := token.Claims.(*claimsJWT)

	if !ok {
		err := errors.New(fmt.Sprintf("token claims are not of type %T", claims))
		return nil, fmt.Errorf("could not cast token claims: %s", err.Error())
	}

	return claims.UserId, nil
}
