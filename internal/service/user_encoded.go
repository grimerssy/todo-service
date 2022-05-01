package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/grimerssy/todo-service/internal/core"
	"github.com/grimerssy/todo-service/internal/repository"
	"github.com/grimerssy/todo-service/pkg/auth"
	"github.com/grimerssy/todo-service/pkg/encoding"
	"github.com/grimerssy/todo-service/pkg/hashing"
)

type UserEncoded struct {
	hasher        hashing.Hasher
	encoder       encoding.Encoder
	authenticator auth.Authenticator
	repository    repository.UserRepository
}

func NewUserEncoded(hasher hashing.Hasher, encoder encoding.Encoder, authenticator auth.Authenticator,
	repository repository.UserRepository) *UserEncoded {

	return &UserEncoded{
		hasher:        hasher,
		encoder:       encoder,
		authenticator: authenticator,
		repository:    repository,
	}
}

func (s *UserEncoded) SignUp(ctx context.Context, userReq core.UserRequest) error {
	user, err := s.requestToUser(userReq)
	if err != nil {
		return fmt.Errorf("could not convert request to user: %s", err.Error())
	}

	if user.Password, err = s.hasher.GenerateHash(user.Password); err != nil {
		return err

	}

	if err := s.repository.Create(ctx, user); err != nil {
		return fmt.Errorf("could not create user: %s", err.Error())
	}

	return nil
}

func (s *UserEncoded) SignIn(ctx context.Context, userReq core.UserRequest) (string, error) {
	user, err := s.repository.GetCredentialsByUsername(ctx, userReq.Username)
	if err != nil {
		return "", ErrUserNotFound
	}

	if match := s.hasher.CompareHashAndPassword(user.Password, userReq.Password); !match {
		return "", errors.New("invalid password")
	}

	id, err := s.encoder.EncodeID(user.ID)
	if err != nil {
		return "", fmt.Errorf("could not encode user id: %s", err.Error())
	}

	token, err := s.authenticator.GenerateToken(id)
	if err != nil {
		return "", fmt.Errorf("could not generate token: %s", err.Error())
	}

	return token, nil
}

func (s *UserEncoded) GetID(ctx context.Context, accessToken string) (any, error) {
	return s.authenticator.ParseToken(accessToken)
}

func (*UserEncoded) requestToUser(su core.UserRequest) (core.User, error) {
	var user core.User

	if len(su.FirstName) == 0 {
		return user, errors.New("empty first name")
	}
	if len(su.LastName) == 0 {
		return user, errors.New("empty last name")
	}
	if len(su.Email) == 0 {
		return user, errors.New("empty email")
	}
	if len(su.Username) == 0 {
		return user, errors.New("empty username")
	}
	if len(su.Password) == 0 {
		return user, errors.New("empty password")
	}

	user = core.User{
		FirstName: su.FirstName,
		LastName:  su.LastName,
		Email:     su.Email,
		Username:  su.Username,
		Password:  su.Password,
	}

	return user, nil
}
