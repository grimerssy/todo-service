package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/grimerssy/todo-service/pkg/service"
)

type Handler struct {
	services service.Service
}

func NewHandler(services service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up")
		auth.POST("/sign-in")
	}

	api := router.Group("/api")
	{
		todos := api.Group("/todos")
		{
			todos.POST("/")
			todos.GET("/:id")
			todos.GET("/pending")
			todos.GET("/")
			todos.PUT("/:id")
			todos.PATCH("/:id")
			todos.DELETE("/:id")
		}
	}

	return router
}
