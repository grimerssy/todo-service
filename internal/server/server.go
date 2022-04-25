package server

import (
	"context"
	"net/http"
)

type ConfigServer struct {
	Http ConfigHttp
}

type ConfigHttp struct {
	Port            string
	RequestSeconds  uint
	ShutdownSeconds uint
}

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(cfg ConfigServer, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:    ":" + cfg.Http.Port,
		Handler: handler,
	}
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context, onShutdown ...func() error) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return err
	}

	for _, f := range onShutdown {
		if err := f(); err != nil {
			return err
		}
	}

	return nil
}
