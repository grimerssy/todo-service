package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/grimerssy/todo-service/internal/core"
	"github.com/grimerssy/todo-service/internal/service"
	"github.com/grimerssy/todo-service/pkg/logging"
)

type TodoGin struct {
	logger         logging.Logger
	todoService    service.TodoService
	requestTimeout time.Duration
}

func NewTodoGin(cfg ConfigGin, logger logging.Logger, todoService service.TodoService) *TodoGin {
	return &TodoGin{
		logger:         logger,
		todoService:    todoService,
		requestTimeout: cfg.RequestSeconds * time.Second,
	}
}

func (h *TodoGin) create(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), h.requestTimeout)
	defer cancel()

	userID, ok := c.Get(userIDKey)
	if !ok {
		err := errors.New("could not get user id")
		h.logger.Log(logging.ErrorLevel, err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	var todoReq core.TodoRequest
	if err := c.BindJSON(&todoReq); err != nil {
		message := "could not bind json"
		h.logger.LogFieldsf(logging.WarnLevel, logging.Fields{
			"user_id": userID,
		}, "could not create todo: %s: %s", message, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": message})
		return
	}

	if err := h.todoService.Create(ctx, userID, todoReq); err != nil {
		message := "could not create todo"
		h.logger.LogFieldsf(logging.ErrorLevel, logging.Fields{
			"user_id": userID,
		}, "%s: %s", message, err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": message})
		return
	}

	h.logger.LogFields(logging.InfoLevel, logging.Fields{
		"user_id": userID,
	}, "created todo")
	c.Status(http.StatusCreated)
}

func (h *TodoGin) getByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), h.requestTimeout)
	defer cancel()

	userID, ok := c.Get(userIDKey)
	if !ok {
		err := errors.New("could not get user id")
		h.logger.Log(logging.ErrorLevel, err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	todoID := c.Param(todoIDKey)

	todoRes, err := h.todoService.GetByID(ctx, userID, todoID)
	switch err {
	case nil:
		h.logger.LogFields(logging.InfoLevel, logging.Fields{
			"user_id": userID,
			"todo_id": todoID,
		}, "got todo by id")
		c.JSON(http.StatusOK, todoRes)
		return
	case service.ErrTodoNotFound:
		h.logger.LogFieldsf(logging.WarnLevel, logging.Fields{
			"user_id": userID,
			"todo_id": todoID,
		}, "could not get todo by id: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	default:
		message := "could not get todo by id"
		h.logger.LogFieldsf(logging.ErrorLevel, logging.Fields{
			"user_id": userID,
			"todo_id": todoID,
		}, "%s: %s", message, err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": message})
	}
}

func (h *TodoGin) getPending(c *gin.Context) {
	const completed = false

	ctx, cancel := context.WithTimeout(context.Background(), h.requestTimeout)
	defer cancel()

	userID, ok := c.Get(userIDKey)
	if !ok {
		err := errors.New("could not get user id")
		h.logger.Log(logging.ErrorLevel, err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	todosRes, err := h.todoService.GetByCompletion(ctx, userID, completed)
	if err != nil {
		message := "could not get pending todos"
		h.logger.LogFieldsf(logging.ErrorLevel, logging.Fields{
			"user_id": userID,
		}, "%s: %s", message, err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": message})
		return
	}

	h.logger.LogFields(logging.InfoLevel, logging.Fields{
		"user_id": userID,
	}, "got pending todos")
	c.JSON(http.StatusOK, todosRes)
}

func (h *TodoGin) getAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), h.requestTimeout)
	defer cancel()

	userID, ok := c.Get(userIDKey)
	if !ok {
		err := errors.New("could not get user id")
		h.logger.Log(logging.ErrorLevel, err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	todosRes, err := h.todoService.GetAll(ctx, userID)
	if err != nil {
		message := "could not get all todos"
		h.logger.LogFieldsf(logging.ErrorLevel, logging.Fields{
			"user_id": userID,
		}, "%s: %s", message, err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": message})
		return
	}

	h.logger.LogFields(logging.InfoLevel, logging.Fields{
		"user_id": userID,
	}, "got all todos")
	c.JSON(http.StatusOK, todosRes)
}

func (h *TodoGin) updateByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), h.requestTimeout)
	defer cancel()

	userID, ok := c.Get(userIDKey)
	if !ok {
		err := errors.New("could not get user id")
		h.logger.Log(logging.ErrorLevel, err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	todoID := c.Param(todoIDKey)

	var todoReq core.TodoRequest
	if err := c.BindJSON(&todoReq); err != nil {
		message := "could not bind json"
		h.logger.LogFieldsf(logging.WarnLevel, logging.Fields{
			"user_id": userID,
			"todo_id": todoID,
		}, "could not update todo by id: %s: %s", message, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": message})
		return
	}

	err := h.todoService.UpdateByID(ctx, userID, todoID, todoReq)

	switch err {
	case nil:
		h.logger.LogFields(logging.InfoLevel, logging.Fields{
			"user_id": userID,
			"todo_id": todoID,
		}, "updated todo by id")
		c.Status(http.StatusNoContent)
		return
	case service.ErrTodoNotFound:
		h.logger.LogFieldsf(logging.WarnLevel, logging.Fields{
			"user_id": userID,
			"todo_id": todoID,
		}, "could not update todo by id: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	default:
		message := "could not update todo by id"
		h.logger.LogFieldsf(logging.ErrorLevel, logging.Fields{
			"user_id": userID,
			"todo_id": todoID,
		}, "%s: %s", message, err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": message})
	}
}

func (h *TodoGin) patchByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), h.requestTimeout)
	defer cancel()

	userID, ok := c.Get(userIDKey)
	if !ok {
		err := errors.New("could not get user id")
		h.logger.Log(logging.ErrorLevel, err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	todoID := c.Param(todoIDKey)

	var todoReq core.TodoRequest
	if err := c.BindJSON(&todoReq); err != nil {
		message := "could not bind json"
		h.logger.LogFieldsf(logging.WarnLevel, logging.Fields{
			"user_id": userID,
			"todo_id": todoID,
		}, "could not patch todo: %s: %s", message, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": message})
		return
	}

	err := h.todoService.PatchByID(ctx, userID, todoID, todoReq)

	switch err {
	case nil:
		h.logger.LogFields(logging.InfoLevel, logging.Fields{
			"user_id": userID,
			"todo_id": todoID,
		}, "patched todo by id")
		c.Status(http.StatusNoContent)
		return
	case service.ErrTodoNotFound:
		h.logger.LogFieldsf(logging.WarnLevel, logging.Fields{
			"user_id": userID,
			"todo_id": todoID,
		}, "could not patch todo by id: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	default:
		message := "could not patch todo by id"
		h.logger.LogFieldsf(logging.ErrorLevel, logging.Fields{
			"user_id": userID,
			"todo_id": todoID,
		}, "%s: %s", message, err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": message})
	}
}

func (h *TodoGin) deleteByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), h.requestTimeout)
	defer cancel()

	userID, ok := c.Get(userIDKey)
	if !ok {
		err := errors.New("could not get user id")
		h.logger.Log(logging.ErrorLevel, err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	todoID := c.Param(todoIDKey)

	err := h.todoService.DeleteByID(ctx, userID, todoID)

	switch err {
	case nil:
		h.logger.LogFields(logging.InfoLevel, logging.Fields{
			"user_id": userID,
			"todo_id": todoID,
		}, "deleted todo by id")
		c.Status(http.StatusNoContent)
		return
	case service.ErrTodoNotFound:
		h.logger.LogFieldsf(logging.WarnLevel, logging.Fields{
			"user_id": userID,
			"todo_id": todoID,
		}, "could not delete todo by id: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	default:
		message := "could not delete todo by id"
		h.logger.LogFieldsf(logging.ErrorLevel, logging.Fields{
			"user_id": userID,
			"todo_id": todoID,
		}, "%s: %s", message, err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": message})
	}
}

func (h *TodoGin) deleteCompleted(c *gin.Context) {
	const completed = true

	ctx, cancel := context.WithTimeout(context.Background(), h.requestTimeout)
	defer cancel()

	userID, ok := c.Get(userIDKey)
	if !ok {
		err := errors.New("could not get user id")
		h.logger.Log(logging.ErrorLevel, err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	if err := h.todoService.DeleteByCompletion(ctx, userID, completed); err != nil {
		message := "could not delete completed todos"
		h.logger.LogFieldsf(logging.ErrorLevel, logging.Fields{
			"user_id": userID,
		}, "%s: %s", message, err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": message})
		return
	}

	h.logger.LogFields(logging.InfoLevel, logging.Fields{
		"user_id": userID,
	}, "deleted completed todos")
	c.Status(http.StatusNoContent)
}
