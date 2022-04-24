package server

import (
	"context"
	"net/http"
	"time"
)

type ConfigHttp struct {
	Port         string
	HeaderBytes  int
	ReadSeconds  int
	WriteSeconds int
}

type HttpServer struct {
	httpServer *http.Server
}

func (s *HttpServer) Run(cfg ConfigHttp, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + cfg.Port,
		Handler:        handler,
		MaxHeaderBytes: cfg.HeaderBytes,
		ReadTimeout:    time.Duration(cfg.ReadSeconds) * time.Second,
		WriteTimeout:   time.Duration(cfg.WriteSeconds) * time.Second,
	}
	return s.httpServer.ListenAndServe()
}

func (s *HttpServer) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
