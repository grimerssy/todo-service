package handler

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/grimerssy/todo-service/pkg/service"
	"github.com/sirupsen/logrus"
)

const (
	authorizationHeader = "Authorization"
	userIDKey           = "user"
)

type MiddlewareGin struct {
	logger         logrus.FieldLogger
	authService    service.AuthService
	requestTimeout time.Duration
}

func NewMiddlewareGin(logger logrus.FieldLogger, authService service.AuthService,
	requestTimeout time.Duration) *MiddlewareGin {

	return &MiddlewareGin{
		logger:         logger,
		authService:    authService,
		requestTimeout: requestTimeout,
	}
}

func (h *MiddlewareGin) authorize(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	ctx := context.TODO()

	if len(header) == 0 {
		message := "could not authorize: empty authorization header"
		h.logger.Errorf("%s", message)
		c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"error": message})
		return
	}

	headerParts := strings.Split(header, " ")

	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		message := "could not authorize: invalid authorization token"
		h.logger.Errorf("%s", message)
		c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"error": message})
		return
	}

	token := headerParts[1]
	userID, err := h.authService.ParseToken(ctx, token)
	if err != nil {
		message := "could not authorize: internal error"
		h.logger.Errorf("%s: %s", message, err.Error())
		c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"error": message})
		return
	}

	c.Set(userIDKey, userID)
}
