package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/grimerssy/todo-service/internal/service"
	"github.com/grimerssy/todo-service/pkg/logging"
)

type MiddlewareGin struct {
	logger         logging.Logger
	userService    service.UserService
	requestTimeout time.Duration
}

func NewMiddlewareGin(cfg ConfigGin, logger logging.Logger, userService service.UserService) *MiddlewareGin {
	return &MiddlewareGin{
		logger:         logger,
		userService:    userService,
		requestTimeout: cfg.RequestSeconds * time.Second,
	}
}

func (h *MiddlewareGin) authorize(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), h.requestTimeout)
	defer cancel()

	header := c.GetHeader(authorizationHeader)

	if len(header) == 0 {
		err := errors.New("could not authorize: empty authorization header")
		h.logger.Log(logging.ErrorLevel, err.Error())
		c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	headerParts := strings.Split(header, " ")

	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		err := errors.New("could not authorize: invalid authorization token")
		h.logger.Log(logging.ErrorLevel, err.Error())
		c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	token := headerParts[1]
	userID, err := h.userService.GetID(ctx, token)
	if err != nil {
		message := "could not authorize"
		h.logger.Logf(logging.ErrorLevel, "%s: %s", message, err.Error())
		c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"error": message})
		return
	}

	c.Set(userIDKey, userID)
}
