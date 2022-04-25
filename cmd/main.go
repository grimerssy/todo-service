package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/grimerssy/todo-service/internal/config"
	"github.com/grimerssy/todo-service/internal/server"
	_ "github.com/lib/pq"
)

const (
	environment = "dev"
)

func main() {
	cfg := config.GetConfig(environment)
	logger := config.GetLogger(cfg.LogFormatting, environment)

	repositories, closeDB := config.GetRepositories(cfg)
	services := config.GetServices(cfg, repositories)
	handlers := config.GetGinHandlers(logger, services)

	srv := new(server.Server)

	go func() {
		if err := srv.Run(cfg.Server, handlers.InitRoutes()); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("an error occured while running http server: %s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if err := srv.Shutdown(context.TODO(), closeDB); err != nil {
		logger.Fatalf("an error occured while shutting down the server: %s", err.Error())
	}
}
