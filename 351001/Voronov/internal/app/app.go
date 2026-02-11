package app

import (
	"Voronov/internal/config"
	"Voronov/pkg/postgres"
	"context"
	"go.uber.org/zap"
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
	
	return nil
}
