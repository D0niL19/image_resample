package server

import (
	"context"
	"imageResample/internal/config"
	"imageResample/internal/service"
	adapter "imageResample/internal/transport/http"
	"log/slog"
	"net/http"
)

type Server struct {
	cfg        *config.Config
	log        *slog.Logger
	httpServer *http.Server
}

func NewServer(cfg *config.Config, storage service.ImageStorage, log *slog.Logger) *Server {
	resampler := service.NewResampler(storage, cfg.ImageWidth, cfg.ImageHeight, log)

	handler := adapter.NewResizeHandler(resampler, log)

	httpServer := &http.Server{
		Addr:    cfg.Address,
		Handler: handler,
	}

	return &Server{
		cfg:        cfg,
		log:        log,
		httpServer: httpServer,
	}
}

func (s *Server) Run() error {
	s.log.Info("Starting HTTP server", slog.String("address", s.cfg.Address))
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.log.Info("Shutting down server...")
	return s.httpServer.Shutdown(ctx)
}
