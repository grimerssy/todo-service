package main

import (
	"log"

	"github.com/grimerssy/todo-service/internal/server"
	"github.com/grimerssy/todo-service/pkg/handler"
)

func main() {
	handlers := new(handler.Handler).InitRoutes()
	srv := new(server.Server)
	if err := srv.Run("8000", handlers); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}
}
