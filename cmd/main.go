package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	handlers := config.GetGinHandlers(cfg, logger, services)

	srv := server.NewServer(cfg.Server, handlers.InitRoutes())

	quit := make(chan os.Signal, 1)

	go func() {
		if err := srv.Run(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("an error occured while running http server: %s", err.Error())
		}
	}()

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownTimeout := cfg.ShutdownSeconds * time.Second
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx, closeDB); err != nil {
		logger.Fatalf("an error occured while shutting down the server: %s", err.Error())
	}
}
