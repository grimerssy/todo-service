package service

import (
	"context"
	"errors"

	"github.com/grimerssy/todo-service/internal/core"
	"github.com/grimerssy/todo-service/pkg/repository"
)

type UserEncoded struct {
	hasher     Hasher
	encoder    Encoder
	repository repository.User
}

func NewUserEncoded(hasher Hasher, encoder Encoder, repository repository.User) *UserEncoded {
	return &UserEncoded{
		hasher:     hasher,
		encoder:    encoder,
		repository: repository,
	}
}

func (s *UserEncoded) Create(ctx context.Context, userSU core.UserSignUp) error {
	user, err := signUpToUser(userSU)
	if err != nil {
		return err
	}

	if user.Password, err = s.hasher.Hash(ctx, user.Password); err != nil {
		return err
	}

	return s.repository.Create(ctx, user)
}

func (s *UserEncoded) GetUserId(ctx context.Context, userSI core.UserSignIn) (interface{}, error) {
	cred, err := s.repository.GetCredentialsByUsername(ctx, userSI.Username)
	if err != nil {
		return nil, err
	}

	if match := s.hasher.CompareHashAndPassword(ctx, cred.Password, userSI.Password); !match {
		return nil, errors.New("invalid password")
	}

	id, err := s.encoder.Encode(ctx, cred.ID)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func signUpToUser(su core.UserSignUp) (core.User, error) {
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
