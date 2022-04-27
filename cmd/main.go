package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/grimerssy/todo-service/internal/config"
	"github.com/grimerssy/todo-service/internal/server"
	"github.com/grimerssy/todo-service/internal/wiring"
	"github.com/grimerssy/todo-service/pkg/logging"
	_ "github.com/lib/pq"
)

const (
	environment = "dev"
)

func main() {
	cfg := config.NewConfig(environment, logging.DefaultLogrus())
	logger := logging.NewLogrus(cfg.Logrus)

	repositories, closeDB := wiring.GetRepositories(cfg, logger)
	services := wiring.GetServices(cfg, logger, repositories)
	handlers := wiring.GetGinHandlers(cfg, logger, services)

	srv := server.NewServer(cfg.Server, handlers.InitRoutes())

	quit := make(chan os.Signal, 1)

	go func() {
		if err := srv.Run(); err != nil && err != http.ErrServerClosed {
			logger.Logf(logging.FatalLevel, "could not run the server: %s", err.Error())
		}
	}()
	logger.Log(logging.InfoLevel, "starting the server")

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log(logging.InfoLevel, "shutting the server down")
	if err := srv.Shutdown(closeDB); err != nil {
		logger.Logf(logging.FatalLevel, "could not shutdown the server: %s", err.Error())
	}
}
