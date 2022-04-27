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

	res := make(chan error, 1)

	go func() {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			res <- err
		}

		for _, f := range onShutdown {
			if err := f(); err != nil {
				res <- err
			}
		}

		res <- nil
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-res:
		return err
	}
}
