package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/grimerssy/todo-service/internal/core"
	"github.com/grimerssy/todo-service/internal/service"
	"github.com/grimerssy/todo-service/pkg/logging"
)

type AuthGin struct {
	logger         logging.Logger
	userService    service.UserService
	requestTimeout time.Duration
}

func NewAuthGin(cfg ConfigGin, logger logging.Logger, userService service.UserService) *AuthGin {
	return &AuthGin{
		logger:         logger,
		userService:    userService,
		requestTimeout: cfg.RequestSeconds * time.Second,
	}
}

func (h *AuthGin) signUp(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), h.requestTimeout)
	defer cancel()

	var userReq core.UserRequest
	if err := c.BindJSON(&userReq); err != nil {
		message := "could not bind json"
		h.logger.Logf(logging.ErrorLevel, "user could not sign up: %s: %s", message, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": message})
		return
	}

	if err := h.userService.SignUp(ctx, userReq); err != nil {
		message := "could not sign up"
		h.logger.Logf(logging.ErrorLevel, "%s: %s", message, err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": message})
		return
	}

	h.logger.LogFields(logging.InfoLevel, logging.Fields{
		"username": userReq.Username,
	}, "new user has signed up")
	c.Status(http.StatusCreated)
}

func (h *AuthGin) signIn(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), h.requestTimeout)
	defer cancel()

	var userReq core.UserRequest
	if err := c.BindJSON(&userReq); err != nil {
		message := "could not bind json"
		h.logger.Logf(logging.ErrorLevel, "user could not sign in: %s: %s", message, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": message})
		return
	}

	token, err := h.userService.SignIn(ctx, userReq)
	switch err {
	case nil:
		h.logger.LogFields(logging.InfoLevel, logging.Fields{
			"username": userReq.Username,
		}, "user has signed in")
		c.JSON(http.StatusOK, map[string]string{"token": token})
		return

	case service.ErrUserNotFound:
		h.logger.LogFieldsf(logging.ErrorLevel, logging.Fields{
			"username": userReq.Username,
		}, "user could not sign in: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		return

	default:
		message := "could not sign in"
		h.logger.LogFields(logging.ErrorLevel, logging.Fields{
			"username": userReq.Username,
		}, "%s: %s", message, err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": message})
	}
}
