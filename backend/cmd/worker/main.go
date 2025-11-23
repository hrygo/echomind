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
	"go.uber.org/zap"
)

// ZapLoggerAdapter adapts zap.Logger to asynq.Logger interface
type ZapLoggerAdapter struct {
	logger *zap.Logger
}

func (l *ZapLoggerAdapter) Debug(args ...interface{}) {
	l.logger.Sugar().Debug(args...)
}

func (l *ZapLoggerAdapter) Info(args ...interface{}) {
	l.logger.Sugar().Info(args...)
}

func (l *ZapLoggerAdapter) Warn(args ...interface{}) {
	l.logger.Sugar().Warn(args...)
}

func (l *ZapLoggerAdapter) Error(args ...interface{}) {
	l.logger.Sugar().Error(args...)
}

func (l *ZapLoggerAdapter) Fatal(args ...interface{}) {
	l.logger.Sugar().Fatal(args...)
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
			Logger:      &ZapLoggerAdapter{logger: container.Logger},
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
			container.Logger.Fatal("Worker failed", zap.Error(err))
		}
	}
}
