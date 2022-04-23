package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/grimerssy/todo-service/internal/core"
)

func (h *Handler) signUp(c *gin.Context) {
	var userReq core.UserRequest
	ctx := context.TODO()

	if err := c.BindJSON(&userReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if err := h.services.UserService.Create(ctx, userReq); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *Handler) signIn(c *gin.Context) {
	var userReq core.UserRequest
	ctx := context.TODO()

	if err := c.BindJSON(&userReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	token, err := h.services.AuthenticationService.GenerateToken(ctx, userReq)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]string{"token": token})
}
