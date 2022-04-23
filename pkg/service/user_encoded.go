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
	repository repository.UserRepository
}

func NewUserEncoded(hasher Hasher, encoder Encoder, repository repository.UserRepository) *UserEncoded {
	return &UserEncoded{
		hasher:     hasher,
		encoder:    encoder,
		repository: repository,
	}
}

func (s *UserEncoded) Create(ctx context.Context, userReq core.UserRequest) error {
	user, err := s.requestToUser(userReq)
	if err != nil {
		return err
	}

	if user.Password, err = s.hasher.Hash(ctx, user.Password); err != nil {
		return err
	}

	return s.repository.Create(ctx, user)
}

func (s *UserEncoded) GetUserId(ctx context.Context, userReq core.UserRequest) (interface{}, error) {
	cred, err := s.repository.GetCredentialsByUsername(ctx, userReq.Username)
	if err != nil {
		return nil, err
	}

	if match := s.hasher.CompareHashAndPassword(ctx, cred.Password, userReq.Password); !match {
		return nil, errors.New("invalid password")
	}

	id, err := s.encoder.Encode(ctx, cred.ID)
	if err != nil {
		return nil, err
	}

	return id, nil
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
