package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/grimerssy/todo-service/internal/core"
	"github.com/grimerssy/todo-service/pkg/service"
	"github.com/sirupsen/logrus"
)

type TodoGin struct {
	logger         logrus.FieldLogger
	todoService    service.TodoService
	requestTimeout time.Duration
}

func NewTodoGin(logger logrus.FieldLogger, todoService service.TodoService,
	requestTimeout time.Duration) *TodoGin {

	return &TodoGin{
		logger:         logger,
		todoService:    todoService,
		requestTimeout: requestTimeout,
	}
}

func (h *TodoGin) create(c *gin.Context) {
	ctx := context.TODO()

	userID, ok := c.Get(userIDKey)
	if !ok {
		message := "could not get user id"
		h.logger.Errorf("%s", message)
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": message})
		return
	}

	var todoReq core.TodoRequest
	if err := c.BindJSON(&todoReq); err != nil {
		message := "could not bind json"
		h.logger.Errorf("%s: %s", message, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": message})
		return
	}

	if err := h.todoService.Create(ctx, userID, todoReq); err != nil {
		message := "could not create user"
		h.logger.Errorf("%s: %s", message, err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": message})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *TodoGin) getByID(c *gin.Context) {
	ctx := context.TODO()

	userID, ok := c.Get(userIDKey)
	if !ok {
		message := "could not get user id"
		h.logger.Errorf("%s", message)
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": message})
		return
	}

	todoID := c.Param("id")

	todoRes, err := h.todoService.GetByID(ctx, userID, todoID)
	if err != nil {
		message := "could not get todo by id"
		h.logger.Errorf("%s: %s", message, err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": message})
		return
	}

	c.JSON(http.StatusOK, todoRes)
}

func (h *TodoGin) getPending(c *gin.Context) {
	const completed = false
	ctx := context.TODO()

	userID, ok := c.Get(userIDKey)
	if !ok {
		message := "could not get user id"
		h.logger.Errorf("%s", message)
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": message})
		return
	}

	todosRes, err := h.todoService.GetByCompletion(ctx, userID, completed)
	if err != nil {
		message := "could not get todos by completion"
		h.logger.Errorf("%s: %s", message, err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": message})
		return
	}

	c.JSON(http.StatusOK, todosRes)
}

func (h *TodoGin) getAll(c *gin.Context) {
	ctx := context.TODO()

	userID, ok := c.Get(userIDKey)
	if !ok {
		message := "could not get user id"
		h.logger.Errorf("%s", message)
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": message})
		return
	}

	todosRes, err := h.todoService.GetAll(ctx, userID)
	if err != nil {
		message := "could not get all todos"
		h.logger.Errorf("%s: %s", message, err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": message})
		return
	}

	c.JSON(http.StatusOK, todosRes)
}

func (h *TodoGin) updateByID(c *gin.Context) {
	ctx := context.TODO()

	userID, ok := c.Get(userIDKey)
	if !ok {
		message := "could not get user id"
		h.logger.Errorf("%s", message)
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": message})
		return
	}

	todoID := c.Param("id")

	var todoReq core.TodoRequest
	if err := c.BindJSON(&todoReq); err != nil {
		message := "could not bind json"
		h.logger.Errorf("%s: %s", message, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": message})
		return
	}

	if err := h.todoService.UpdateByID(ctx, userID, todoID, todoReq); err != nil {
		message := "could not update todo by id"
		h.logger.Errorf("%s: %s", message, err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": message})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *TodoGin) patchByID(c *gin.Context) {
	ctx := context.TODO()

	userID, ok := c.Get(userIDKey)
	if !ok {
		message := "could not get user id"
		h.logger.Errorf("%s", message)
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": message})
		return
	}

	todoID := c.Param("id")

	var todoReq core.TodoRequest
	if err := c.BindJSON(&todoReq); err != nil {
		message := "could not bind json"
		h.logger.Errorf("%s: %s", message, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": message})
		return
	}

	if err := h.todoService.PatchByID(ctx, userID, todoID, todoReq); err != nil {
		message := "could not patch todo by id"
		h.logger.Errorf("%s: %s", message, err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": message})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *TodoGin) deleteByID(c *gin.Context) {
	ctx := context.TODO()

	userID, ok := c.Get(userIDKey)
	if !ok {
		message := "could not get user id"
		h.logger.Errorf("%s", message)
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": message})
		return
	}

	todoID := c.Param("id")

	if err := h.todoService.DeleteByID(ctx, userID, todoID); err != nil {
		message := "could not delete todo by id"
		h.logger.Errorf("%s: %s", message, err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": message})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *TodoGin) deleteCompleted(c *gin.Context) {
	const completed = true
	ctx := context.TODO()

	userID, ok := c.Get(userIDKey)
	if !ok {
		message := "could not get user id"
		h.logger.Errorf("%s", message)
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": message})
		return
	}

	if err := h.todoService.DeleteByCompletion(ctx, userID, completed); err != nil {
		message := "could not delete todo by completion"
		h.logger.Errorf("%s: %s", message, err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": message})
		return
	}

	c.Status(http.StatusNoContent)
}
