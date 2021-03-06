package server

import (
	"context"
	"net/http"
	"time"
)

type ConfigServer struct {
	ShutdownSeconds time.Duration
	Http            struct {
		Port string
	}
}

type Server struct {
	shutdownTimeout time.Duration
	httpServer      *http.Server
}

func NewServer(cfg ConfigServer, handler http.Handler) *Server {
	return &Server{
		shutdownTimeout: cfg.ShutdownSeconds * time.Second,
		httpServer: &http.Server{
			Addr:    ":" + cfg.Http.Port,
			Handler: handler,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(onShutdown ...func() error) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return err
	}

	for _, fn := range onShutdown {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}
