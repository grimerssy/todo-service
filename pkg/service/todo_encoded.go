package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/grimerssy/todo-service/internal/core"
	"github.com/grimerssy/todo-service/pkg/repository"
)

type TodoEncoded struct {
	userEncoder Encoder
	todoEncoder Encoder
	repository  repository.TodoRepository
}

func NewTodoEncoded(userEncoder, todoEncoder Encoder, repository repository.TodoRepository) *TodoEncoded {
	return &TodoEncoded{
		userEncoder: userEncoder,
		todoEncoder: todoEncoder,
		repository:  repository,
	}
}

func (s *TodoEncoded) Create(ctx context.Context, userID interface{}, todoReq core.TodoRequest) error {
	uintUserID, err := s.userEncoder.Decode(ctx, userID)
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

	return nil
}

func (s *TodoEncoded) GetByID(ctx context.Context, userID interface{}, todoID interface{}) (core.TodoResponse, error) {
	var response core.TodoResponse

	uintUserID, err := s.userEncoder.Decode(ctx, userID)
	if err != nil {
		return response, fmt.Errorf("could not decode user id: %s", err.Error())
	}

	uintTodoID, err := s.todoEncoder.Decode(ctx, todoID)
	if err != nil {
		return response, fmt.Errorf("could not decode todo id: %s", err.Error())
	}

	todo, err := s.repository.GetByID(ctx, uintUserID, uintTodoID)
	if err != nil {
		return response, fmt.Errorf("could not get todo: %s", err.Error())
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

	uintUserID, err := s.userEncoder.Decode(ctx, userID)
	if err != nil {
		return responses, fmt.Errorf("could not decode user id: %s", err.Error())
	}

	todos, err := s.repository.GetByCompletion(ctx, uintUserID, completed)
	if err != nil {
		return responses, fmt.Errorf("could not get todos: %s", err.Error())
	}

	responses = make([]core.TodoResponse, len(todos))

	for i, todo := range todos {
		todoID, err := s.todoEncoder.Encode(ctx, todo.ID)
		if err != nil {
			return responses, fmt.Errorf("could not encode todo id: %s", err.Error())
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

	uintUserID, err := s.userEncoder.Decode(ctx, userID)
	if err != nil {
		return responses, fmt.Errorf("could not decode user id: %s", err.Error())
	}

	todos, err := s.repository.GetAll(ctx, uintUserID)
	if err != nil {
		return responses, fmt.Errorf("could not get todos: %s", err.Error())
	}

	responses = make([]core.TodoResponse, len(todos))

	for i, todo := range todos {
		todoID, err := s.todoEncoder.Encode(ctx, todo.ID)
		if err != nil {
			return responses, fmt.Errorf("could not encode todo id: %s", err.Error())
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
	uintUserID, err := s.userEncoder.Decode(ctx, userID)
	if err != nil {
		return fmt.Errorf("could not decode user id: %s", err.Error())
	}

	todo, err := s.requestToTodo(todoReq)
	if err != nil {
		return fmt.Errorf("could not convert request to todo: %s", err.Error())
	}

	uintTodoID, err := s.todoEncoder.Decode(ctx, todoID)
	if err != nil {
		return fmt.Errorf("could not decode todo id: %s", err.Error())
	}

	if err := s.repository.UpdateByID(ctx, uintUserID, uintTodoID, todo); err != nil {
		return fmt.Errorf("could not update todo: %s", err.Error())
	}

	return nil
}

func (s *TodoEncoded) PatchByID(ctx context.Context, userID interface{}, todoID interface{}, todoReq core.TodoRequest) error {
	uintUserID, err := s.userEncoder.Decode(ctx, userID)
	if err != nil {
		return fmt.Errorf("could not decode user id: %s", err.Error())
	}

	todo := core.Todo{
		Title:       todoReq.Title,
		Description: todoReq.Description,
		Completed:   todoReq.Completed,
	}

	uintTodoID, err := s.todoEncoder.Decode(ctx, todoID)
	if err != nil {
		return fmt.Errorf("could not decode todo id: %s", err.Error())
	}

	if err := s.repository.PatchByID(ctx, uintUserID, uintTodoID, todo); err != nil {
		return fmt.Errorf("could not patch todo: %s", err.Error())
	}

	return nil
}

func (s *TodoEncoded) DeleteByID(ctx context.Context, userID interface{}, todoID interface{}) error {
	uintUserID, err := s.userEncoder.Decode(ctx, userID)
	if err != nil {
		return fmt.Errorf("could not decode user id: %s", err.Error())
	}

	uintTodoID, err := s.todoEncoder.Decode(ctx, todoID)
	if err != nil {
		return fmt.Errorf("could not decode todo id: %s", err.Error())
	}

	if err := s.repository.DeleteByID(ctx, uintUserID, uintTodoID); err != nil {
		return fmt.Errorf("could not delete todo: %s", err.Error())
	}

	return nil
}

func (s *TodoEncoded) DeleteByCompletion(ctx context.Context, userID interface{}, completed bool) error {
	uintUserID, err := s.userEncoder.Decode(ctx, userID)
	if err != nil {
		return fmt.Errorf("could not decode user id: %s", err.Error())
	}

	if err := s.repository.DeleteByCompletion(ctx, uintUserID, completed); err != nil {
		return fmt.Errorf("could not delete todos: %s", err.Error())
	}

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
