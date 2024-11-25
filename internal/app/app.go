package app

import (
	"context"
	"errors"
	"imageResample/internal/config"
	"imageResample/internal/server"
	"imageResample/internal/storage"
	"imageResample/pkg/logger"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() {
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env)

	dirStorage := storage.NewDirectoryStorage(cfg.PathOrig, cfg.PathRes)

	srv := server.NewServer(cfg, dirStorage, log)

	go func() {
		if err := srv.Run(); err != nil && !errors.Is(err, context.Canceled) {
			log.Error("Server failed", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Failed to shutdown server gracefully", slog.Any("error", err))
	}
	log.Info("Server exited gracefully")
}
