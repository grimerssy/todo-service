package handler

import (
	"github.com/gin-gonic/gin"
)

type HandlersGin struct {
	auth       AuthGin
	middleware MiddlewareGin
	todo       TodoGin
}

func (h *HandlersGin) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.auth.signUp)
		auth.POST("/sign-in", h.auth.signIn)
	}

	api := router.Group("/api", h.middleware.authorize)
	{
		todos := api.Group("/todos")
		{
			todos.POST("/", h.todo.create)
			todos.GET("/:id", h.todo.getByID)
			todos.GET("/pending", h.todo.getPending)
			todos.GET("/", h.todo.getAll)
			todos.PUT("/:id", h.todo.updateByID)
			todos.PATCH("/:id", h.todo.patchByID)
			todos.DELETE("/:id", h.todo.deleteByID)
			todos.DELETE("/completed", h.todo.deleteCompleted)
		}
	}

	return router
}
