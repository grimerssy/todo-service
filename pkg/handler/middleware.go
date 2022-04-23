package handler

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userIDKey           = "user"
)

func (h *Handler) authorize(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	ctx := context.TODO()

	if len(header) == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"error": "empty authorization header"})
		return
	}

	headerParts := strings.Split(header, " ")

	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"error": "invalid authorization header"})
		return
	}

	token := headerParts[1]
	userID, err := h.services.AuthenticationService.ParseToken(ctx, token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	c.Set(userIDKey, userID)
}
