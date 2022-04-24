package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/grimerssy/todo-service/cmd/config"
	"github.com/grimerssy/todo-service/internal/server"
	_ "github.com/lib/pq"
)

const (
	environment = "dev"
)

func main() {
	cfg := config.GetConfig(environment)
	logger := config.GetLogger(cfg.LogFormatting, environment)

	db, repositories := config.GetDbAndRepositories(cfg)
	services := config.GetServices(cfg, repositories)
	handlers := config.GetGinHandlers(logger, services)

	srv := new(server.HttpServer)

	go func() {
		if err := srv.Run(cfg.Http, handlers.InitRoutes()); err != nil {
			log.Fatalf("an error occured while running http server: %s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if err := srv.Shutdown(context.TODO()); err != nil {
		log.Fatalf("an error occured while shutting down the server: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		log.Fatalf("an error occured while closing db connection: %s", err.Error())
	}
}
