package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"subscriptions/internal/config"
	docsHandler "subscriptions/internal/handler/docs"
	subscriptionHandler "subscriptions/internal/handler/subscription"
	postgres "subscriptions/internal/infrastructure/postgresql"
	"subscriptions/internal/middleware"
	subscriptionRepository "subscriptions/internal/repository/subscription"
	subscriptionService "subscriptions/internal/service/subscription"
)

func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	logger := newLogger(cfg.Log.Level)
	pool, err := postgres.New(context.Background(), cfg.Database.DSN())
	if err != nil {
		return fmt.Errorf("connect postgres: %w", err)
	}
	defer pool.Close()

	repo := subscriptionRepository.NewRepository(pool)
	svc := subscriptionService.NewService(repo)
	h := subscriptionHandler.NewHandler(svc, logger)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("ok")) })
	docsHandler.RegisterRoutes(mux)
	subscriptionHandler.RegisterRoutes(mux, h)

	handler := middleware.Chain(
		mux,
		middleware.RequestID,
		middleware.Logging(logger),
		middleware.Recovery(logger),
		middleware.Timeout(30*time.Second),
	)
	server := &http.Server{
		Addr:         cfg.Server.Addr(),
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	go func() {
		logger.Info("server started", "addr", cfg.Server.Addr())
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "err", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return server.Shutdown(ctx)
}

func newLogger(level string) *slog.Logger {
	logLevel := slog.LevelInfo
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	}
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
}
