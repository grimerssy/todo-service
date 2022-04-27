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
	res := make(chan error, 1)

	go func() {
		user, err := s.requestToUser(userReq)
		if err != nil {
			res <- fmt.Errorf("could not convert request to user: %s", err.Error())
			return
		}

		if user.Password, err = s.hasher.GenerateHash(ctx, user.Password); err != nil {
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

func (s *UserEncoded) SignIn(ctx context.Context, userReq core.UserRequest) (string, error) {
	res := make(chan func() (string, error), 1)

	go func() {
		user, err := s.repository.GetCredentialsByUsername(ctx, userReq.Username)
		if err != nil {
			res <- func() (string, error) {
				return "", ErrUserNotFound
			}
			return
		}

		if match := s.hasher.CompareHashAndPassword(ctx, user.Password, userReq.Password); !match {
			res <- func() (string, error) {
				return "", errors.New("invalid password")
			}
			return
		}

		id, err := s.encoder.EncodeID(ctx, user.ID)
		if err != nil {
			res <- func() (string, error) {
				return "", fmt.Errorf("could not encode user id: %s", err.Error())
			}
			return
		}

		token, err := s.authenticator.GenerateToken(ctx, id)
		if err != nil {
			res <- func() (string, error) {
				return "", fmt.Errorf("could not generate token: %s", err.Error())
			}
		}

		res <- func() (string, error) {
			return token, nil
		}
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case f := <-res:
		return f()
	}
}

func (s *UserEncoded) GetID(ctx context.Context, token string) (interface{}, error) {
	res := make(chan func() (interface{}, error), 1)

	go func() {
		res <- func() (interface{}, error) {
			return s.authenticator.ParseToken(ctx, token)
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
