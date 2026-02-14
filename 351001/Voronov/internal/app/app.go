package app

import (
	"Voronov/internal/config"
	"Voronov/pkg/postgres"
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func Run(ctx context.Context, logger *zap.Logger) error {
	cfg, err := config.New()
	if err != nil {
		logger.Fatal("load config: %w", zap.Error(err))
	}

	db, err := postgres.New(cfg.Postgres)
	if err != nil {
		logger.Fatal("connect to database: %w", zap.Error(err))
	}
	defer db.Close(ctx)

	mux := http.NewServeMux()

	logger.Info("starting server", zap.String("host", "localhost"), zap.String("port", cfg.HTTPport))

	// написать хендлеры

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.HTTPport),
		Handler: mux,
	}

	go func() {
		logger.Info("starting server", zap.String("port", cfg.HTTPport))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {

			logger.Fatal("listen and serve failed", zap.Error(err))
		}
	}()
	<-ctx.Done()

	logger.Info("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	logger.Info("server successfully shutdown...")

	return nil
}
