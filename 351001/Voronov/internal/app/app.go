package app

import (
	"Voronov/internal/config"
	"Voronov/internal/repository"
	"Voronov/internal/service"
	"Voronov/internal/transport/handler"
	"Voronov/pkg/postgres"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

func Run(ctx context.Context, logger *zap.Logger) error {
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("validate config: %w", err)
	}

	// Run migrations via database/sql (goose requirement)
	db, err := sql.Open("postgres", cfg.GooseDBString)
	if err != nil {
		return fmt.Errorf("open db for migrations: %w", err)
	}
	// Закроем это соединение, когда миграции отработают
	defer db.Close()

	// 2. Проверяем, что база доступна
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping db: %w", err)
	}

	// 3. САМОЕ ВАЖНОЕ: Создаем схему ПЕРЕД настройкой goose
	_, err = db.ExecContext(ctx, "CREATE SCHEMA IF NOT EXISTS distcomp;")
	if err != nil {
		return fmt.Errorf("create schema distcomp: %w", err)
	}

	// 4. Настраиваем goose (явно указываем схему для служебной таблицы)
	goose.SetTableName("distcomp.schema_migrations")
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set goose dialect: %w", err)
	}

	// 5. Запускаем миграции (убедись, что папка migrations лежит рядом с main.go или укажи правильный путь)
	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}
	logger.Info("migrations applied successfully")

	// Create pgxpool for repositories
	pool, err := postgres.NewPool(ctx, &postgres.Config{
		Host:     cfg.PostgresHost,
		Port:     cfg.PostgresPort,
		Username: cfg.PostgresUser,
		Password: cfg.PostgresPass,
		Database: cfg.PostgresDB,
	})
	if err != nil {
		return fmt.Errorf("create pool: %w", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("ping pool: %w", err)
	}

	// Wire repositories
	userRepo := repository.NewUserRepository(pool)
	issueRepo := repository.NewIssueRepository(pool)
	labelRepo := repository.NewLabelRepository(pool)
	reactionRepo := repository.NewReactionRepository(pool)

	// Wire services
	mapper := service.NewMapper()
	userService := service.NewUserService(userRepo, mapper)
	issueService := service.NewIssueService(issueRepo, userRepo, labelRepo, reactionRepo, mapper)
	labelService := service.NewLabelService(labelRepo, mapper)
	reactionService := service.NewReactionService(reactionRepo, issueRepo, mapper)

	// Wire handler
	h := handler.NewHandler(userService, issueService, labelService, reactionService)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.HTTPport),
		Handler: mux,
	}

	go func() {
		logger.Info("server listening", zap.String("addr", server.Addr))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("listen and serve failed", zap.Error(err))
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	logger.Info("server shutdown complete")
	return nil
}
