package server

import (
	"context"
	"net/http"
)

type ConfigServer struct {
	Http ConfigHttp
}

type ConfigHttp struct {
	Port string
}

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg ConfigServer, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    ":" + cfg.Http.Port,
			Handler: handler,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context, onShutdown ...func() error) error {
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
