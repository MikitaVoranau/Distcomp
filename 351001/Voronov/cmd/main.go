package main

import (
	"Voronov/internal/app"
	"Voronov/pkg/logger"
	"context"
	"go.uber.org/zap"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger, err := logger.New()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	if err := app.Run(ctx, logger); err != nil {
		logger.Fatal("Error starting app", zap.Error(err))
	}
}
