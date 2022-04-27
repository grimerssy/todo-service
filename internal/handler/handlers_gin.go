package handler

import (
	"time"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userIDKey           = "user_id"
	todoIDKey           = "todo_id"
)

type ConfigGin struct {
	RequestSeconds time.Duration
}

type HandlersGin struct {
	Auth       *AuthGin
	Middleware *MiddlewareGin
	Todo       *TodoGin
}

func (h *HandlersGin) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.Auth.signUp)
		auth.POST("/sign-in", h.Auth.signIn)
	}

	api := router.Group("/api", h.Middleware.authorize)
	{
		todos := api.Group("/todos")
		{
			todos.POST("/", h.Todo.create)
			todos.GET("/:"+todoIDKey, h.Todo.getByID)
			todos.GET("/pending", h.Todo.getPending)
			todos.GET("/", h.Todo.getAll)
			todos.PUT("/:"+todoIDKey, h.Todo.updateByID)
			todos.PATCH("/:"+todoIDKey, h.Todo.patchByID)
			todos.DELETE("/:"+todoIDKey, h.Todo.deleteByID)
			todos.DELETE("/completed", h.Todo.deleteCompleted)
		}
	}

	return router
}
