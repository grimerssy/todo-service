package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/grimerssy/todo-service/internal/core"
	"github.com/grimerssy/todo-service/internal/repository"
	"github.com/grimerssy/todo-service/pkg/cache"
	"github.com/grimerssy/todo-service/pkg/encoding"
)

type TodoEncoded struct {
	cache       cache.Cache
	userEncoder encoding.Encoder
	todoEncoder encoding.Encoder
	repository  repository.TodoRepository
}

func NewTodoEncoded(cache cache.Cache, userEncoder, todoEncoder encoding.Encoder,
	repository repository.TodoRepository) *TodoEncoded {

	return &TodoEncoded{
		cache:       cache,
		userEncoder: userEncoder,
		todoEncoder: todoEncoder,
		repository:  repository,
	}
}

func (s *TodoEncoded) Create(ctx context.Context, userID interface{}, todoReq core.TodoRequest) error {
	res := make(chan error, 1)

	go func() {
		uintUserID, err := s.userEncoder.DecodeID(ctx, userID)
		if err != nil {
			res <- fmt.Errorf("could not decode user id: %s", err.Error())
			return
		}

		todo, err := s.requestToTodo(todoReq)
		if err != nil {
			res <- fmt.Errorf("could not convert request to todo: %s", err.Error())
			return
		}

		if err := s.repository.Create(ctx, uintUserID, todo); err != nil {
			res <- fmt.Errorf("could not create todo: %s", err.Error())
			return
		}

		s.invalidateUserCache(uintUserID, nil)

		res <- nil
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-res:
		return err
	}
}

func (s *TodoEncoded) GetByID(ctx context.Context, userID interface{}, todoID interface{}) (core.TodoResponse, error) {
	res := make(chan func() (core.TodoResponse, error), 1)

	go func() {
		uintUserID, err := s.userEncoder.DecodeID(ctx, userID)
		if err != nil {
			res <- func() (core.TodoResponse, error) {
				return core.TodoResponse{}, fmt.Errorf("could not decode user id: %s", err.Error())
			}
			return
		}

		uintTodoID, err := s.todoEncoder.DecodeID(ctx, todoID)
		if err != nil {
			res <- func() (core.TodoResponse, error) {
				return core.TodoResponse{}, fmt.Errorf("could not decode todo id: %s", err.Error())
			}
			return
		}

		cacheKey := cache.TodoCacheKey{
			UserID: uintUserID,
			Args:   uintTodoID,
		}

		if cached := s.cache.GetValue(cacheKey); cached != nil {
			res <- func() (core.TodoResponse, error) {
				return cached.(core.TodoResponse), nil
			}
			return
		}

		todo, err := s.repository.GetByID(ctx, uintUserID, uintTodoID)
		if err != nil {
			res <- func() (core.TodoResponse, error) {
				return core.TodoResponse{}, ErrTodoNotFound
			}
			return
		}

		response := core.TodoResponse{
			ID:          todoID,
			Title:       todo.Title,
			Description: todo.Description,
			Completed:   todo.Completed,
		}

		s.cache.SetValue(cacheKey, response)

		res <- func() (core.TodoResponse, error) {
			return response, nil
		}
	}()

	select {
	case <-ctx.Done():
		return core.TodoResponse{}, ctx.Err()
	case fn := <-res:
		return fn()
	}
}

func (s *TodoEncoded) GetByCompletion(ctx context.Context, userID interface{}, completed bool) ([]core.TodoResponse, error) {
	res := make(chan func() ([]core.TodoResponse, error), 1)

	go func() {
		uintUserID, err := s.userEncoder.DecodeID(ctx, userID)
		if err != nil {
			res <- func() ([]core.TodoResponse, error) {
				return nil, fmt.Errorf("could not decode user id: %s", err.Error())
			}
			return
		}

		cacheKey := cache.TodoCacheKey{
			UserID: uintUserID,
			Args:   completed,
		}

		if cached := s.cache.GetValue(cacheKey); cached != nil {
			res <- func() ([]core.TodoResponse, error) {
				return cached.([]core.TodoResponse), nil
			}
			return
		}

		todos, err := s.repository.GetByCompletion(ctx, uintUserID, completed)
		if err != nil {
			res <- func() ([]core.TodoResponse, error) {
				return nil, fmt.Errorf("could not get todos: %s", err.Error())
			}
			return
		}

		responses := make([]core.TodoResponse, len(todos))

		for i, todo := range todos {
			todoID, err := s.todoEncoder.EncodeID(ctx, todo.ID)
			if err != nil {
				res <- func() ([]core.TodoResponse, error) {
					return nil, fmt.Errorf("could not encode todo id: %s", err.Error())
				}
				return
			}

			responses[i] = core.TodoResponse{
				ID:          todoID,
				Title:       todo.Title,
				Description: todo.Description,
				Completed:   todo.Completed,
			}
		}

		s.cache.SetValue(cacheKey, responses)

		res <- func() ([]core.TodoResponse, error) {
			return responses, nil
		}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case fn := <-res:
		return fn()
	}
}

func (s *TodoEncoded) GetAll(ctx context.Context, userID interface{}) ([]core.TodoResponse, error) {
	res := make(chan func() ([]core.TodoResponse, error), 1)

	go func() {
		uintUserID, err := s.userEncoder.DecodeID(ctx, userID)
		if err != nil {
			res <- func() ([]core.TodoResponse, error) {
				return nil, fmt.Errorf("could not decode user id: %s", err.Error())
			}
			return
		}

		cacheKey := cache.TodoCacheKey{
			UserID: uintUserID,
			Args:   nil,
		}

		if cached := s.cache.GetValue(cacheKey); cached != nil {
			res <- func() ([]core.TodoResponse, error) {
				return cached.([]core.TodoResponse), nil
			}
			return
		}

		todos, err := s.repository.GetAll(ctx, uintUserID)
		if err != nil {
			res <- func() ([]core.TodoResponse, error) {
				return nil, fmt.Errorf("could not get todos: %s", err.Error())
			}
			return
		}

		responses := make([]core.TodoResponse, len(todos))

		for i, todo := range todos {
			todoID, err := s.todoEncoder.EncodeID(ctx, todo.ID)
			if err != nil {
				res <- func() ([]core.TodoResponse, error) {
					return nil, fmt.Errorf("could not encode todo id: %s", err.Error())
				}
				return
			}

			responses[i] = core.TodoResponse{
				ID:          todoID,
				Title:       todo.Title,
				Description: todo.Description,
				Completed:   todo.Completed,
			}
		}

		s.cache.SetValue(cacheKey, responses)

		res <- func() ([]core.TodoResponse, error) {
			return responses, nil
		}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case fn := <-res:
		return fn()
	}
}

func (s *TodoEncoded) UpdateByID(ctx context.Context, userID interface{}, todoID interface{}, todoReq core.TodoRequest) error {
	res := make(chan error, 1)

	go func() {
		uintUserID, err := s.userEncoder.DecodeID(ctx, userID)
		if err != nil {
			res <- fmt.Errorf("could not decode user id: %s", err.Error())
			return
		}

		todo, err := s.requestToTodo(todoReq)
		if err != nil {
			res <- fmt.Errorf("could not convert request to todo: %s", err.Error())
			return
		}

		uintTodoID, err := s.todoEncoder.DecodeID(ctx, todoID)
		if err != nil {
			res <- fmt.Errorf("could not decode todo id: %s", err.Error())
			return
		}

		id, err := s.repository.UpdateByID(ctx, uintUserID, uintTodoID, todo)
		if err != nil {
			res <- fmt.Errorf("could not update todo: %s", err.Error())
			return
		}

		s.invalidateUserCache(uintUserID, []uint{id})

		res <- nil
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-res:
		return err
	}
}

func (s *TodoEncoded) PatchByID(ctx context.Context, userID interface{}, todoID interface{}, todoReq core.TodoRequest) error {
	res := make(chan error, 1)

	go func() {
		uintUserID, err := s.userEncoder.DecodeID(ctx, userID)
		if err != nil {
			res <- fmt.Errorf("could not decode user id: %s", err.Error())
			return
		}

		todo := core.Todo{
			Title:       todoReq.Title,
			Description: todoReq.Description,
			Completed:   todoReq.Completed,
		}

		uintTodoID, err := s.todoEncoder.DecodeID(ctx, todoID)
		if err != nil {
			res <- fmt.Errorf("could not decode todo id: %s", err.Error())
			return
		}

		id, err := s.repository.PatchByID(ctx, uintUserID, uintTodoID, todo)
		if err != nil {
			res <- fmt.Errorf("could not patch todo: %s", err.Error())
			return
		}

		s.invalidateUserCache(uintUserID, []uint{id})

		res <- nil
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-res:
		return err
	}
}

func (s *TodoEncoded) DeleteByID(ctx context.Context, userID interface{}, todoID interface{}) error {
	res := make(chan error, 1)

	go func() {
		uintUserID, err := s.userEncoder.DecodeID(ctx, userID)
		if err != nil {
			res <- fmt.Errorf("could not decode user id: %s", err.Error())
			return
		}

		uintTodoID, err := s.todoEncoder.DecodeID(ctx, todoID)
		if err != nil {
			res <- fmt.Errorf("could not decode todo id: %s", err.Error())
			return
		}

		id, err := s.repository.DeleteByID(ctx, uintUserID, uintTodoID)
		if err != nil {
			res <- fmt.Errorf("could not delete todo: %s", err.Error())
			return
		}

		s.invalidateUserCache(uintUserID, []uint{id})

		res <- nil
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-res:
		return err
	}
}

func (s *TodoEncoded) DeleteByCompletion(ctx context.Context, userID interface{}, completed bool) error {
	res := make(chan error, 1)

	go func() {
		uintUserID, err := s.userEncoder.DecodeID(ctx, userID)
		if err != nil {
			res <- fmt.Errorf("could not decode user id: %s", err.Error())
			return
		}

		ids, err := s.repository.DeleteByCompletion(ctx, uintUserID, completed)
		if err != nil {
			res <- fmt.Errorf("could not delete todos: %s", err.Error())
			return
		}

		s.invalidateUserCache(uintUserID, ids)

		res <- nil
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-res:
		return err
	}
}

func (*TodoEncoded) requestToTodo(req core.TodoRequest) (core.Todo, error) {
	var todo core.Todo

	if len(req.Title) == 0 {
		return todo, errors.New("empty title")
	}

	todo = core.Todo{
		Title:       req.Title,
		Description: req.Description,
		Completed:   false,
	}

	return todo, nil
}

func (s *TodoEncoded) invalidateUserCache(userID uint, todoIDs []uint) {
	cacheKeys := []cache.TodoCacheKey{
		{
			UserID: userID,
			Args:   true,
		},
		{
			UserID: userID,
			Args:   false,
		},
		{
			UserID: userID,
			Args:   nil,
		},
	}
	for _, todoID := range todoIDs {
		cacheKey := cache.TodoCacheKey{
			UserID: userID,
			Args:   todoID,
		}
		cacheKeys = append(cacheKeys, cacheKey)
	}

	for _, key := range cacheKeys {
		s.cache.RemoveValue(key)
	}
}
