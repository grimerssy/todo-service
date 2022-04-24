package handler

import (
	"github.com/gin-gonic/gin"
)

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
			todos.GET("/:id", h.Todo.getByID)
			todos.GET("/pending", h.Todo.getPending)
			todos.GET("/", h.Todo.getAll)
			todos.PUT("/:id", h.Todo.updateByID)
			todos.PATCH("/:id", h.Todo.patchByID)
			todos.DELETE("/:id", h.Todo.deleteByID)
			todos.DELETE("/completed", h.Todo.deleteCompleted)
		}
	}

	return router
}
