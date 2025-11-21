package main

import (
	"context"
	"log"
    "strings"

	"github.com/hrygo/echomind/internal/service"
	"github.com/hrygo/echomind/internal/tasks"
	"github.com/hibiken/asynq"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Initialize Viper
	vip := viper.New()
	vip.SetConfigFile("backend/configs/config.yaml")
	vip.AddConfigPath(".")
    
    // Enable Environment Variable Overrides
    vip.AutomaticEnv()
    vip.SetEnvPrefix("ECHOMIND")
    vip.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
    
	if err := vip.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	// Logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Database
	dsn := vip.GetString("database.dsn")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// AI Service
	aiProvider, err := service.AIProviderFactory(vip)
	if err != nil {
		log.Fatalf("Failed to create AI provider: %v", err)
	}
	summarizer := service.NewSummaryService(aiProvider)

	// Redis & Asynq
	redisAddr := vip.GetString("redis.addr")
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{Concurrency: 10},
	)

	// Mux
	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeEmailAnalyze, func(ctx context.Context, t *asynq.Task) error {
		return tasks.HandleEmailAnalyzeTask(ctx, t, db, summarizer)
	})

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}