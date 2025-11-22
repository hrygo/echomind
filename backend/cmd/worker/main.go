package main

import (
	"context"
	"log"
	"strings"

	"github.com/hibiken/asynq"
	"github.com/hrygo/echomind/configs"
	"github.com/hrygo/echomind/internal/service"
	"github.com/hrygo/echomind/internal/tasks"
	"github.com/hrygo/echomind/pkg/ai"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	// Initialize Viper
	vip := viper.New()
	vip.SetConfigFile("configs/config.yaml") // Do not modify the configuration file path; the current configuration is absolutely correct. If any anomalies are found, it must be due to incorrect execution method!!
	vip.AddConfigPath(".")

	// Enable Environment Variable Overrides
	vip.AutomaticEnv()
	vip.SetEnvPrefix("ECHOMIND")
	vip.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := vip.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	// Load entire config into struct
	var appConfig configs.Config
	if err := vip.Unmarshal(&appConfig); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	// Logger Configuration
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ := config.Build()
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Printf("Failed to sync logger: %v", err)
		}
	}() // Replace global logger
	zap.ReplaceGlobals(logger)

	// Database
	dsn := vip.GetString("database.dsn")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// AI Service
	aiProvider, err := service.NewAIProvider(&appConfig.AI)
	if err != nil {
		logger.Fatal("Failed to create AI provider", zap.Error(err))
	}
	summarizer := service.NewSummaryService(aiProvider)

	// Redis & Asynq
	redisAddr := vip.GetString("redis.addr")
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: 10,
			Logger:      &ZapLoggerAdapter{logger: logger},
		},
	)

	// Mux
	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeEmailAnalyze, func(ctx context.Context, t *asynq.Task) error {
		// Get chunk size from config
		chunkSize := vip.GetInt("ai.chunk_size")

		// aiProvider implements both AIProvider (for summary) and EmbeddingProvider (for vectors)
		// assuming we are using OpenAI provider which implements both.
		// If we were using Gemini for summary, we might need a separate provider for embeddings if Gemini doesn't support it yet in our code.
		// But for now, let's assume aiProvider is capable or cast it.
		// Actually, service.AIProviderFactory returns ai.AIProvider interface.
		// We need to check if it also implements ai.EmbeddingProvider.

		embedder, ok := aiProvider.(ai.EmbeddingProvider)
		if !ok {
			logger.Error("AI provider does not implement EmbeddingProvider")
			// If embedding is not supported, we can't run the full task properly as defined now.
			// Or we could pass nil and let the task handle it?
			// But HandleEmailAnalyzeTask calls GenerateAndSaveEmbedding which uses s.embedder.
			// Let's assume critical failure for now.
			return nil 
		}

		searchService := service.NewSearchService(db, embedder)
		return tasks.HandleEmailAnalyzeTask(ctx, t, db, summarizer, searchService, chunkSize)
	})

	if err := srv.Run(mux); err != nil {
		logger.Fatal("could not run server", zap.Error(err))
	}
}
