package service

import (
	"context"
	"errors"

	"github.com/grimerssy/todo-service/internal/core"
	"github.com/grimerssy/todo-service/pkg/repository"
)

type TodoEncoded struct {
	encoder    Encoder
	repository repository.TodoRepository
}

func NewTodoEncoded(encoder Encoder, repository repository.TodoRepository) *TodoEncoded {
	return &TodoEncoded{
		encoder:    encoder,
		repository: repository,
	}
}

func (s *TodoEncoded) Create(ctx context.Context, userID interface{}, todoReq core.TodoRequest) error {
	todo, err := s.requestToTodo(todoReq)
	if err != nil {
		return err
	}

	uintUserID, err := s.encoder.Decode(ctx, userID)
	if err != nil {
		return err
	}

	return s.repository.Create(ctx, uintUserID, todo)
}

func (s *TodoEncoded) GetByID(ctx context.Context, userID interface{}, todoID interface{}) (core.TodoResponse, error) {
	var response core.TodoResponse

	uintUserID, err := s.encoder.Decode(ctx, userID)
	if err != nil {
		return response, err
	}

	uintTodoID, err := s.encoder.Decode(ctx, todoID)
	if err != nil {
		return response, err
	}

	todo, err := s.repository.GetByID(ctx, uintUserID, uintTodoID)
	if err != nil {
		return response, err
	}

	response = core.TodoResponse{
		Id:          todoID,
		Title:       todo.Title,
		Description: todo.Description,
		Completed:   todo.Completed,
	}

	return response, nil
}

func (s *TodoEncoded) GetByCompletion(ctx context.Context, userID interface{}, completed bool) ([]core.TodoResponse, error) {
	var responses []core.TodoResponse

	uintUserID, err := s.encoder.Decode(ctx, userID)
	if err != nil {
		return responses, err
	}

	todos, err := s.repository.GetByCompletion(ctx, uintUserID, completed)
	if err != nil {
		return responses, err
	}

	responses = make([]core.TodoResponse, len(todos))

	for i, todo := range todos {
		todoID, err := s.encoder.Encode(ctx, todo.ID)
		if err != nil {
			return responses, err
		}

		responses[i] = core.TodoResponse{
			Id:          todoID,
			Title:       todo.Title,
			Description: todo.Description,
			Completed:   todo.Completed,
		}
	}

	return responses, nil
}

func (s *TodoEncoded) GetAll(ctx context.Context, userID interface{}) ([]core.TodoResponse, error) {
	var responses []core.TodoResponse

	uintUserID, err := s.encoder.Decode(ctx, userID)
	if err != nil {
		return responses, err
	}

	todos, err := s.repository.GetAll(ctx, uintUserID)
	if err != nil {
		return responses, err
	}

	responses = make([]core.TodoResponse, len(todos))

	for i, todo := range todos {
		todoID, err := s.encoder.Encode(ctx, todo.ID)
		if err != nil {
			return responses, err
		}

		responses[i] = core.TodoResponse{
			Id:          todoID,
			Title:       todo.Title,
			Description: todo.Description,
			Completed:   todo.Completed,
		}
	}

	return responses, nil
}

func (s *TodoEncoded) UpdateByID(ctx context.Context, userID interface{}, todoID interface{}, todoReq core.TodoRequest) error {
	todo, err := s.requestToTodo(todoReq)
	if err != nil {
		return err
	}

	uintUserID, err := s.encoder.Decode(ctx, userID)
	if err != nil {
		return err
	}

	uintTodoID, err := s.encoder.Decode(ctx, todoID)
	if err != nil {
		return err
	}

	return s.repository.UpdateByID(ctx, uintUserID, uintTodoID, todo)
}

func (s *TodoEncoded) PatchByID(ctx context.Context, userID interface{}, todoID interface{}, todoReq core.TodoRequest) error {
	todo := core.Todo{
		Title:       todoReq.Title,
		Description: todoReq.Description,
		Completed:   todoReq.Completed,
	}

	uintUserID, err := s.encoder.Decode(ctx, userID)
	if err != nil {
		return err
	}

	uintTodoID, err := s.encoder.Decode(ctx, todoID)
	if err != nil {
		return err
	}

	return s.repository.PatchByID(ctx, uintUserID, uintTodoID, todo)
}

func (s *TodoEncoded) DeleteByID(ctx context.Context, userID interface{}, todoID interface{}) error {
	uintUserID, err := s.encoder.Decode(ctx, userID)
	if err != nil {
		return err
	}

	uintTodoID, err := s.encoder.Decode(ctx, todoID)
	if err != nil {
		return err
	}

	return s.repository.DeleteByID(ctx, uintUserID, uintTodoID)
}

func (s *TodoEncoded) DeleteByCompletion(ctx context.Context, userID interface{}, completed bool) error {
	uintUserID, err := s.encoder.Decode(ctx, userID)
	if err != nil {
		return err
	}

	return s.repository.DeleteByCompletion(ctx, uintUserID, completed)
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
