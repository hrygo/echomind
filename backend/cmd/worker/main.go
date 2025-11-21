package main

import (
	"context"
	"log"
	"strings"

	"github.com/hibiken/asynq"
	"github.com/hrygo/echomind/internal/service"
	"github.com/hrygo/echomind/internal/tasks"
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

	// Logger Configuration
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ := config.Build()
	                defer func() {
	                        if err := logger.Sync(); err != nil {
	                                log.Printf("Failed to sync logger: %v", err)
	                        }
	                }()	// Replace global logger
	zap.ReplaceGlobals(logger)

	// Database
	dsn := vip.GetString("database.dsn")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// AI Service
	aiProvider, err := service.AIProviderFactory(vip)
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
		return tasks.HandleEmailAnalyzeTask(ctx, t, db, summarizer)
	})

	if err := srv.Run(mux); err != nil {
		logger.Fatal("could not run server", zap.Error(err))
	}
}
