package main

import (
	"fmt"
	"log"
    "strings"
	"time"

	"echomind.com/backend/internal/handler"
	"echomind.com/backend/internal/model"
	"echomind.com/backend/internal/service"
	clientimap "github.com/emersion/go-imap/client"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const Version = "0.1.0"

func main() {
	// Initialize Viper for configuration
	vip := viper.New()
	vip.SetConfigFile("configs/config.yaml") // Use the new config path
	vip.AddConfigPath(".")
    
    // Enable Environment Variable Overrides
    vip.AutomaticEnv()
    vip.SetEnvPrefix("ECHOMIND") // e.g., ECHOMIND_AI_DEEPSEEK_API_KEY will override ai.deepseek.api_key
    vip.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := vip.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	// Initialize Zap logger
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	sugar.Infof("Logger initialized. EchoMind Version: %s", Version)

	// Initialize GORM database
	dsn := vip.GetString("database.dsn")
	if dsn == "" {
		sugar.Fatal("Database DSN not found in config")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		sugar.Fatalf("Failed to connect to database: %v", err)
	}

	// AutoMigrate the Email model
	if err := db.AutoMigrate(&model.Email{}); err != nil {
		sugar.Fatalf("Failed to auto migrate database: %v", err)
	}
	sugar.Infof("Database migration completed")

	// Initialize Asynq Client
	redisAddr := vip.GetString("redis.addr")
	asynqClient := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	defer asynqClient.Close()

	// Initialize Gin router
	r := gin.Default()
	
	// Enable CORS for frontend development
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Dependencies for handlers
	imapClient := &clientimap.Client{}
	defaultFetcher := &service.DefaultFetcher{}

	// Initialize Handlers
	syncHandler := handler.NewSyncHandler(db, imapClient, defaultFetcher, asynqClient)
	emailHandler := handler.NewEmailHandler(db)

	// Register routes
	api := r.Group("/api/v1")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})
		api.POST("/sync", syncHandler.SyncEmails)
		api.GET("/emails", emailHandler.ListEmails)
		api.GET("/emails/:id", emailHandler.GetEmail)
	}

	port := vip.GetString("server.port")
	if port == "" {
		port = "8080"
	}

	sugar.Infof("Starting server on :%s", port)
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		sugar.Fatalf("Failed to run server: %v", err)
	}
}
