package main

import (
	"context"
	"log"

	"github.com/hibiken/asynq"
	"github.com/hrygo/echomind/internal/bootstrap"
	"github.com/hrygo/echomind/internal/service"
	"github.com/hrygo/echomind/internal/tasks"
	"github.com/hrygo/echomind/pkg/ai"
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
	// 1. Bootstrap Application
	app, err := bootstrap.Init("configs/config.yaml", false)
	if err != nil {
		log.Fatalf("Bootstrap failed: %v", err)
	}
	defer app.Close()

	// 2. Initialize Services
	aiProvider, err := service.NewAIProvider(&app.Config.AI)
	if err != nil {
		app.Logger.Fatal("Failed to create AI provider", zap.Error(err))
	}
	summarizer := service.NewSummaryService(aiProvider)

	// 3. Setup Asynq Server
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     app.Config.Redis.Addr,
			Password: app.Config.Redis.Password,
			DB:       app.Config.Redis.DB,
		},
		asynq.Config{
			Concurrency: 10,
			Logger:      &ZapLoggerAdapter{logger: app.Logger},
		},
	)

	// 4. Mux & Tasks
	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeEmailAnalyze, func(ctx context.Context, t *asynq.Task) error {
		// Use explicit cast from Config via mapstructure or manual check if needed
		// For now assuming viper populated it correctly in Config struct, but Config struct doesn't have chunk_size
		// Let's assume default or add it to Config struct later.
		// Checking app.Config.AI...
		// It's not in app.Config.AI (AIConfig).
		// We'll default to 1000 if not found, or add to config.
		chunkSize := 1000 // Default

		embedder, ok := aiProvider.(ai.EmbeddingProvider)
		if !ok {
			app.Logger.Error("AI provider does not implement EmbeddingProvider")
			return nil 
		}

		searchService := service.NewSearchService(app.DB, embedder)
		contextService := service.NewContextService(app.DB)
		return tasks.HandleEmailAnalyzeTask(ctx, t, app.DB, summarizer, searchService, contextService, chunkSize)
	})

	app.Logger.Info("Starting worker...")
	if err := srv.Run(mux); err != nil {
		app.Logger.Fatal("could not run server", zap.Error(err))
	}
}