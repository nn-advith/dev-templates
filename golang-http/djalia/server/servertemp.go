package server

import (
	"context"
	"net/http"
	"time"
)

type Config struct {
	Address      string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type AppServer struct {
	Server http.Server
}

func NewServer(cfg Config, handler http.Handler) (*AppServer, error) {
	return &AppServer{
		Server: http.Server{
			Addr:         cfg.Address,
			Handler:      handler,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
	}, nil
}

func (a *AppServer) Start() error {
	return a.Server.ListenAndServe()
}

func (a *AppServer) Stop(ctx context.Context) error {
	return a.Server.Shutdown(ctx)
}
