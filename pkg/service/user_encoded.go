package service

import (
	"context"
	"errors"
	"fmt"

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
	res := make(chan error, 1)

	go func() {
		user, err := s.requestToUser(userReq)
		if err != nil {
			res <- fmt.Errorf("could not convert request to user: %s", err.Error())
			return
		}

		if user.Password, err = s.hasher.Hash(ctx, user.Password); err != nil {
			res <- err
			return
		}

		if err := s.repository.Create(ctx, user); err != nil {
			res <- fmt.Errorf("could not create user: %s", err.Error())
			return
		}

		res <- nil
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-res:
		return err
	}
}

func (s *UserEncoded) GetUserId(ctx context.Context, userReq core.UserRequest) (interface{}, error) {
	res := make(chan func() (interface{}, error), 1)

	go func() {
		cred, err := s.repository.GetCredentialsByUsername(ctx, userReq.Username)
		if err != nil {
			res <- func() (interface{}, error) {
				return nil, fmt.Errorf("could not get user credentials: %s", err.Error())
			}
			return
		}

		if match := s.hasher.CompareHashAndPassword(ctx, cred.Password, userReq.Password); !match {
			res <- func() (interface{}, error) {
				return nil, errors.New("invalid password")
			}
			return
		}

		id, err := s.encoder.Encode(ctx, cred.ID)
		if err != nil {
			res <- func() (interface{}, error) {
				return nil, fmt.Errorf("could not encode user id: %s", err.Error())
			}
			return
		}

		res <- func() (interface{}, error) {
			return id, nil
		}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case f := <-res:
		return f()
	}
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
