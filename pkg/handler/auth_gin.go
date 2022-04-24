package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/grimerssy/todo-service/internal/core"
	"github.com/grimerssy/todo-service/pkg/service"
)

type AuthGin struct {
	authService service.AuthService
	userService service.UserService
}

func (h *AuthGin) signUp(c *gin.Context) {
	var userReq core.UserRequest
	ctx := context.TODO()

	if err := c.BindJSON(&userReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if err := h.userService.Create(ctx, userReq); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *AuthGin) signIn(c *gin.Context) {
	var userReq core.UserRequest
	ctx := context.TODO()

	if err := c.BindJSON(&userReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	token, err := h.authService.GenerateToken(ctx, userReq)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]string{"token": token})
}
