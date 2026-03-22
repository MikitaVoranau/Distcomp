package app

import (
	"Voronov/internal/config"
	"Voronov/internal/model"
	"Voronov/internal/repository"
	"Voronov/internal/service"
	"Voronov/internal/transport/handler"
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

	if err := cfg.Validate(); err != nil {
		logger.Fatal("validate config: %w", zap.Error(err))
	}

	userRepo := repository.NewInMemoryRepository(
		func(u *model.User) int64 { return u.ID },
		func(u *model.User, id int64) { u.ID = id },
	)
	issueRepo := repository.NewInMemoryRepository(
		func(i *model.Issue) int64 { return i.ID },
		func(i *model.Issue, id int64) { i.ID = id },
	)
	labelRepo := repository.NewInMemoryRepository(
		func(l *model.Label) int64 { return l.ID },
		func(l *model.Label, id int64) { l.ID = id },
	)
	reactionRepo := repository.NewInMemoryRepository(
		func(r *model.Reaction) int64 { return r.ID },
		func(r *model.Reaction, id int64) { r.ID = id },
	)
	issueLabelRepo := repository.NewInMemoryRepository(
		func(il *model.IssueLabel) int64 { return il.IssueID*1000 + il.LabelID },
		func(il *model.IssueLabel, id int64) { il.IssueID = id / 1000; il.LabelID = id % 1000 },
	)

	mapper := service.NewMapper()

	userService := service.NewUserService(userRepo, mapper)
	issueService := service.NewIssueService(issueRepo, userRepo, labelRepo, reactionRepo, issueLabelRepo, mapper)
	labelService := service.NewLabelService(labelRepo, mapper)
	reactionService := service.NewReactionService(reactionRepo, mapper)

	h := handler.NewHandler(userService, issueService, labelService, reactionService)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	logger.Info("starting server", zap.String("host", "localhost"), zap.String("port", cfg.HTTPport))

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.HTTPport),
		Handler: mux,
	}

	go func() {
		logger.Info("server listening", zap.String("port", cfg.HTTPport))
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
