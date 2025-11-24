package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hibiken/asynq"
	"github.com/hrygo/echomind/internal/app"
	"github.com/hrygo/echomind/internal/tasks"
	"github.com/hrygo/echomind/pkg/logger"
)

// LoggerAdapter adapts our logger to asynq.Logger interface
type LoggerAdapter struct {
	logger logger.Logger
}

func (l *LoggerAdapter) Debug(args ...interface{}) {
	l.logger.Debug("Asynq Debug", logger.Any("args", args))
}

func (l *LoggerAdapter) Info(args ...interface{}) {
	l.logger.Info("Asynq Info", logger.Any("args", args))
}

func (l *LoggerAdapter) Warn(args ...interface{}) {
	l.logger.Warn("Asynq Warn", logger.Any("args", args))
}

func (l *LoggerAdapter) Error(args ...interface{}) {
	l.logger.Error("Asynq Error", logger.Any("args", args))
}

func (l *LoggerAdapter) Fatal(args ...interface{}) {
	l.logger.Fatal("Asynq Fatal", logger.Any("args", args))
}

func main() {
	// Parse CLI configuration
	cli := app.ParseCLI()

	// Initialize application container
	container, err := app.NewContainer(cli.ConfigPath, cli.IsProduction)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	defer container.Close()

	// Setup Asynq Server
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     container.Config.Redis.Addr,
			Password: container.Config.Redis.Password,
			DB:       container.Config.Redis.DB,
		},
		asynq.Config{
			Concurrency: container.WorkerConcurrency(),
			Logger:      &LoggerAdapter{logger: container.Logger},
		},
	)

	// Register task handlers
	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeEmailAnalyze, func(ctx context.Context, t *asynq.Task) error {
		return tasks.HandleEmailAnalyzeTask(
			ctx, t,
			container.DB,
			container.Summarizer,
			container.SearchService,
			container.ContextService,
			container.ChunkSize(),
			container.Logger, // Use new logger
		)
	})
	mux.HandleFunc(tasks.TypeEmailSync, func(ctx context.Context, t *asynq.Task) error {
		return tasks.HandleEmailSyncTask(
			ctx, t,
			container.SyncService,
			container.Logger, // Use new logger
		)
	})

	container.Logger.Info("Starting worker...")

	// Run worker in a goroutine
	done := make(chan error, 1)
	go func() {
		done <- srv.Run(mux)
	}()

	// Wait for interrupt signal to gracefully shutdown the worker
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		container.Logger.Info("Shutting down worker...")
		srv.Shutdown()
		container.Logger.Info("Worker stopped gracefully")
	case err := <-done:
		if err != nil {
			container.Logger.Fatal("Worker failed", logger.Error(err))
		}
	}
}
