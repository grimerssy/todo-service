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

type AuthGin struct {
	logger         logrus.FieldLogger
	authService    service.AuthService
	userService    service.UserService
	requestTimeout time.Duration
}

func NewAuthGin(logger logrus.FieldLogger, authService service.AuthService, userService service.UserService,
	requestTimeout time.Duration) *AuthGin {

	return &AuthGin{
		logger:         logger,
		authService:    authService,
		userService:    userService,
		requestTimeout: requestTimeout,
	}
}

func (h *AuthGin) signUp(c *gin.Context) {
	var userReq core.UserRequest
	ctx := context.TODO()

	if err := c.BindJSON(&userReq); err != nil {
		message := "could not bind json"
		h.logger.Errorf("%s: %s", message, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": message})
		return
	}

	if err := h.userService.Create(ctx, userReq); err != nil {
		message := "could not create user"
		h.logger.Errorf(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": message})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *AuthGin) signIn(c *gin.Context) {
	var userReq core.UserRequest
	ctx := context.TODO()

	if err := c.BindJSON(&userReq); err != nil {
		message := "could not bind json"
		h.logger.Errorf("%s: %s", message, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": message})
		return
	}

	token, err := h.authService.GenerateToken(ctx, userReq)
	if err != nil {
		message := "could not sign in"
		h.logger.Errorf("%s: %s", message, err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": message})
		return
	}

	c.JSON(http.StatusOK, map[string]string{"token": token})
}
