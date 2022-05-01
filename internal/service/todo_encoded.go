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

func (s *TodoEncoded) Create(ctx context.Context, userID any, todoReq core.TodoRequest) error {
	uintUserID, err := s.userEncoder.DecodeID(userID)
	if err != nil {
		return fmt.Errorf("could not decode user id: %s", err.Error())
	}

	todo, err := s.requestToTodo(todoReq)
	if err != nil {
		return fmt.Errorf("could not convert request to todo: %s", err.Error())
	}

	if err := s.repository.Create(ctx, uintUserID, todo); err != nil {
		return fmt.Errorf("could not create todo: %s", err.Error())
	}

	s.invalidateUserCache(uintUserID, nil)

	return nil
}

func (s *TodoEncoded) GetByID(ctx context.Context, userID, todoID any) (core.TodoResponse, error) {
	uintUserID, err := s.userEncoder.DecodeID(userID)
	if err != nil {
		return core.TodoResponse{}, fmt.Errorf("could not decode user id: %s", err.Error())
	}

	uintTodoID, err := s.todoEncoder.DecodeID(todoID)
	if err != nil {
		return core.TodoResponse{}, fmt.Errorf("could not decode todo id: %s", err.Error())
	}

	cacheKey := cache.TodoCacheKey{
		UserID: uintUserID,
		Args:   uintTodoID,
	}

	if cached := s.cache.GetValue(cacheKey); cached != nil {
		return cached.(core.TodoResponse), nil
	}

	todo, err := s.repository.GetByID(ctx, uintUserID, uintTodoID)
	if err != nil {
		return core.TodoResponse{}, ErrTodoNotFound
	}

	response := core.TodoResponse{
		ID:          todoID,
		Title:       todo.Title,
		Description: todo.Description,
		Completed:   todo.Completed,
	}

	s.cache.SetValue(cacheKey, response)

	return response, nil
}

func (s *TodoEncoded) GetByCompletion(ctx context.Context, userID any, completed bool) ([]core.TodoResponse, error) {
	uintUserID, err := s.userEncoder.DecodeID(userID)
	if err != nil {
		return nil, fmt.Errorf("could not decode user id: %s", err.Error())
	}

	cacheKey := cache.TodoCacheKey{
		UserID: uintUserID,
		Args:   completed,
	}

	if cached := s.cache.GetValue(cacheKey); cached != nil {
		return cached.([]core.TodoResponse), nil
	}

	todos, err := s.repository.GetByCompletion(ctx, uintUserID, completed)
	if err != nil {
		return nil, fmt.Errorf("could not get todos: %s", err.Error())
	}

	responses := make([]core.TodoResponse, len(todos))

	for i, todo := range todos {
		todoID, err := s.todoEncoder.EncodeID(todo.ID)
		if err != nil {
			return nil, fmt.Errorf("could not encode todo id: %s", err.Error())
		}

		responses[i] = core.TodoResponse{
			ID:          todoID,
			Title:       todo.Title,
			Description: todo.Description,
			Completed:   todo.Completed,
		}
	}

	s.cache.SetValue(cacheKey, responses)

	return responses, nil
}

func (s *TodoEncoded) GetAll(ctx context.Context, userID any) ([]core.TodoResponse, error) {
	uintUserID, err := s.userEncoder.DecodeID(userID)
	if err != nil {
		return nil, fmt.Errorf("could not decode user id: %s", err.Error())
	}

	cacheKey := cache.TodoCacheKey{
		UserID: uintUserID,
		Args:   nil,
	}

	if cached := s.cache.GetValue(cacheKey); cached != nil {
		return cached.([]core.TodoResponse), nil
	}

	todos, err := s.repository.GetAll(ctx, uintUserID)
	if err != nil {
		return nil, fmt.Errorf("could not get todos: %s", err.Error())
	}

	responses := make([]core.TodoResponse, len(todos))

	for i, todo := range todos {
		todoID, err := s.todoEncoder.EncodeID(todo.ID)
		if err != nil {
			return nil, fmt.Errorf("could not encode todo id: %s", err.Error())
		}

		responses[i] = core.TodoResponse{
			ID:          todoID,
			Title:       todo.Title,
			Description: todo.Description,
			Completed:   todo.Completed,
		}
	}

	s.cache.SetValue(cacheKey, responses)

	return responses, nil
}

func (s *TodoEncoded) UpdateByID(ctx context.Context, userID, todoID any, todoReq core.TodoRequest) error {
	uintUserID, err := s.userEncoder.DecodeID(userID)
	if err != nil {
		return fmt.Errorf("could not decode user id: %s", err.Error())
	}

	todo, err := s.requestToTodo(todoReq)
	if err != nil {
		return fmt.Errorf("could not convert request to todo: %s", err.Error())
	}

	uintTodoID, err := s.todoEncoder.DecodeID(todoID)
	if err != nil {
		return fmt.Errorf("could not decode todo id: %s", err.Error())
	}

	id, err := s.repository.UpdateByID(ctx, uintUserID, uintTodoID, todo)
	if err != nil {
		return ErrTodoNotFound
	}

	s.invalidateUserCache(uintUserID, []uint{id})

	return nil
}

func (s *TodoEncoded) PatchByID(ctx context.Context, userID, todoID any, todoReq core.TodoRequest) error {
	uintUserID, err := s.userEncoder.DecodeID(userID)
	if err != nil {
		return fmt.Errorf("could not decode user id: %s", err.Error())
	}

	todo := core.Todo{
		Title:       todoReq.Title,
		Description: todoReq.Description,
		Completed:   todoReq.Completed,
	}

	uintTodoID, err := s.todoEncoder.DecodeID(todoID)
	if err != nil {
		return fmt.Errorf("could not decode todo id: %s", err.Error())
	}

	id, err := s.repository.PatchByID(ctx, uintUserID, uintTodoID, todo)
	if err != nil {
		return ErrTodoNotFound
	}

	s.invalidateUserCache(uintUserID, []uint{id})

	return nil
}

func (s *TodoEncoded) DeleteByID(ctx context.Context, userID, todoID any) error {
	uintUserID, err := s.userEncoder.DecodeID(userID)
	if err != nil {
		return fmt.Errorf("could not decode user id: %s", err.Error())
	}

	uintTodoID, err := s.todoEncoder.DecodeID(todoID)
	if err != nil {
		return fmt.Errorf("could not decode todo id: %s", err.Error())
	}

	id, err := s.repository.DeleteByID(ctx, uintUserID, uintTodoID)
	if err != nil {
		return ErrTodoNotFound
	}

	s.invalidateUserCache(uintUserID, []uint{id})

	return nil
}

func (s *TodoEncoded) DeleteByCompletion(ctx context.Context, userID any, completed bool) error {
	uintUserID, err := s.userEncoder.DecodeID(userID)
	if err != nil {
		return fmt.Errorf("could not decode user id: %s", err.Error())
	}

	ids, err := s.repository.DeleteByCompletion(ctx, uintUserID, completed)
	if err != nil {
		return fmt.Errorf("could not delete todos: %s", err.Error())
	}

	s.invalidateUserCache(uintUserID, ids)

	return nil
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
