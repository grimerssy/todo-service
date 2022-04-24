package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/grimerssy/todo-service/internal/core"
	"github.com/grimerssy/todo-service/pkg/service"
)

type TodoGin struct {
	todoService service.TodoService
}

func (h *TodoGin) create(c *gin.Context) {
	ctx := context.TODO()

	userID, ok := c.Get(userIDKey)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": "could not get user id"})
		return
	}

	var todoReq core.TodoRequest
	if err := c.BindJSON(&todoReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if err := h.todoService.Create(ctx, userID, todoReq); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *TodoGin) getByID(c *gin.Context) {
	ctx := context.TODO()

	userID, ok := c.Get(userIDKey)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": "could not get user id"})
		return
	}

	todoID := c.Param("id")
	if len(todoID) == 0 {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": "invalid todo id"})
		return
	}

	todoRes, err := h.todoService.GetByID(ctx, userID, todoID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, todoRes)
}

func (h *TodoGin) getPending(c *gin.Context) {
	const completed = false
	ctx := context.TODO()

	userID, ok := c.Get(userIDKey)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": "could not get user id"})
		return
	}

	todosRes, err := h.todoService.GetByCompletion(ctx, userID, completed)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, todosRes)
}

func (h *TodoGin) getAll(c *gin.Context) {
	ctx := context.TODO()

	userID, ok := c.Get(userIDKey)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": "could not get user id"})
		return
	}

	todosRes, err := h.todoService.GetAll(ctx, userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, todosRes)
}

func (h *TodoGin) updateByID(c *gin.Context) {
	ctx := context.TODO()

	userID, ok := c.Get(userIDKey)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": "could not get user id"})
		return
	}

	todoID := c.Param("id")
	if len(todoID) == 0 {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": "invalid todo id"})
		return
	}

	var todoReq core.TodoRequest
	if err := c.BindJSON(&todoReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if err := h.todoService.UpdateByID(ctx, userID, todoID, todoReq); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *TodoGin) patchByID(c *gin.Context) {
	ctx := context.TODO()

	userID, ok := c.Get(userIDKey)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": "could not get user id"})
		return
	}

	todoID := c.Param("id")
	if len(todoID) == 0 {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": "invalid todo id"})
		return
	}

	var todoReq core.TodoRequest
	if err := c.BindJSON(&todoReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if err := h.todoService.PatchByID(ctx, userID, todoID, todoReq); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *TodoGin) deleteByID(c *gin.Context) {
	ctx := context.TODO()

	userID, ok := c.Get(userIDKey)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": "could not get user id"})
		return
	}

	todoID := c.Param("id")
	if len(todoID) == 0 {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": "invalid todo id"})
		return
	}

	if err := h.todoService.DeleteByID(ctx, userID, todoID); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *TodoGin) deleteCompleted(c *gin.Context) {
	const completed = true
	ctx := context.TODO()

	userID, ok := c.Get(userIDKey)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": "could not get user id"})
		return
	}

	if err := h.todoService.DeleteByCompletion(ctx, userID, completed); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
