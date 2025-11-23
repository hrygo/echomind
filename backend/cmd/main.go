package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/hrygo/echomind/configs"
	"github.com/hrygo/echomind/internal/handler"
	"github.com/hrygo/echomind/internal/middleware"
	"github.com/hrygo/echomind/internal/model"
	"github.com/hrygo/echomind/internal/service"
	"github.com/hrygo/echomind/pkg/ai"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const Version = "0.9.0"

func main() {
	// Initialize Viper for configuration
	vip := viper.New()
	vip.SetConfigFile("configs/config.yaml") // Do not modify the configuration file path; the current configuration is absolutely correct. If any anomalies are found, it must be due to incorrect execution method!!
	vip.AddConfigPath(".")

	// Enable Environment Variable Overrides
	vip.AutomaticEnv()
	vip.SetEnvPrefix("ECHOMIND") // e.g., ECHOMIND_AI_DEEPSEEK_API_KEY will override ai.deepseek.api_key
	vip.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := vip.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	// Load entire config into struct
	var appConfig configs.Config
	if err := vip.Unmarshal(&appConfig); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	// Initialize Zap logger
	logger, _ := zap.NewProduction()
	sugar := logger.Sugar()
	defer func() {
		if err := logger.Sync(); err != nil {
			sugar.Errorf("Failed to sync logger: %v", err)
		}
	}() // flushes buffer, if any
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

	// Enable pgvector extension
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS vector").Error; err != nil {
		sugar.Fatalf("Failed to create vector extension: %v", err)
	}

	// AutoMigrate models
	if err := db.AutoMigrate(&model.Email{}, &model.User{}, &model.Contact{}, &model.EmailAccount{}, &model.EmailEmbedding{}, &model.Organization{}, &model.OrganizationMember{}, &model.Team{}, &model.TeamMember{}, &model.Task{}, &model.Context{}, &model.EmailContext{}); err != nil {
		sugar.Fatalf("Failed to auto migrate database: %v", err)
	}
	sugar.Infof("Database migration completed")

	// Create HNSW index for vector search optimization
	// Note: vector_cosine_ops is for Cosine distance (default for text-embedding-3-small usually)
	if err := db.Exec("CREATE INDEX IF NOT EXISTS email_embeddings_vector_idx ON email_embeddings USING hnsw (vector vector_cosine_ops)").Error; err != nil {
		sugar.Warnf("Failed to create HNSW index (might be expected if table is empty or index exists): %v", err)
	}

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
	defaultFetcher := &service.DefaultFetcher{}

	// Initialize Services
	aiProvider, err := service.NewAIProvider(&appConfig.AI)
	if err != nil {
		sugar.Fatalf("Failed to create AI provider: %v", err)
	}

	organizationService := service.NewOrganizationService(db)
	userService := service.NewUserService(db, appConfig.Server.JWT, organizationService)
	emailService := service.NewEmailService(db)
	contactService := service.NewContactService(db)
	accountService := service.NewAccountService(db, &appConfig.Security)
	insightService := service.NewInsightService(db)
	aiDraftService := service.NewAIDraftService(aiProvider)
	syncService := service.NewSyncService(db, &service.DefaultIMAPClient{}, defaultFetcher, asynqClient, contactService, accountService, &appConfig, sugar)

	// Run Organization Migration (Ensure existing users have an org)
	if err := organizationService.EnsureAllUsersHaveOrganization(context.Background()); err != nil {
		sugar.Errorf("Failed to migrate organizations: %v", err)
		// We don't fatal here to avoid blocking startup in case of minor issues, but in production this should be monitored
	}

	// Cast aiProvider to EmbeddingProvider for search
	embedder, ok := aiProvider.(ai.EmbeddingProvider)
	if !ok {
		sugar.Fatal("AI provider does not implement EmbeddingProvider")
	}
	searchService := service.NewSearchService(db, embedder)
	chatService := service.NewChatService(aiProvider, searchService)
	taskService := service.NewTaskService(db)
	contextService := service.NewContextService(db)

	// Initialize Handlers
	accountHandler := handler.NewAccountHandler(accountService)
	syncHandler := handler.NewSyncHandler(syncService)
	emailHandler := handler.NewEmailHandler(emailService)
	authHandler := handler.NewAuthHandler(userService)
	insightHandler := handler.NewInsightHandler(insightService)
	aiDraftHandler := handler.NewAIDraftHandler(aiDraftService)
	searchHandler := handler.NewSearchHandler(searchService, sugar)
	healthHandler := handler.NewHealthHandler(db)
	orgHandler := handler.NewOrganizationHandler(organizationService)
	chatHandler := handler.NewChatHandler(chatService)
	taskHandler := handler.NewTaskHandler(taskService)
	contextHandler := handler.NewContextHandler(contextService)

	// Register routes
	api := r.Group("/api/v1")
	{
		api.GET("/health", healthHandler.HealthCheck)
		// Auth Routes
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)

		// Protected routes (require JWT authentication)
		protected := api.Group("/").Use(middleware.AuthMiddleware(appConfig.Server.JWT))
		{
			// Organization Routes
			protected.POST("/orgs", orgHandler.CreateOrganization)
			protected.GET("/orgs", orgHandler.ListOrganizations)
			protected.GET("/orgs/:id", orgHandler.GetOrganization)
			protected.GET("/orgs/:id/members", orgHandler.GetMembers)

			protected.POST("/settings/account", accountHandler.ConnectAndSaveAccount)
			protected.GET("/settings/account", accountHandler.GetAccountStatus)
			protected.POST("/sync", syncHandler.SyncEmails)
			protected.GET("/emails", emailHandler.ListEmails)
			protected.GET("/emails/:id", emailHandler.GetEmail)
			protected.DELETE("/emails/all", emailHandler.DeleteAllEmails)
			protected.GET("/insights/network", insightHandler.GetNetworkGraph)
			protected.POST("/ai/draft", aiDraftHandler.GenerateDraft)
			protected.GET("/search", searchHandler.Search)
			protected.POST("/chat/completions", chatHandler.StreamChat)

			// Task Routes
			protected.POST("/tasks", taskHandler.CreateTask)
			protected.GET("/tasks", taskHandler.ListTasks)
			protected.PATCH("/tasks/:id", taskHandler.UpdateTask)
			protected.PATCH("/tasks/:id/status", taskHandler.UpdateTaskStatus)
			protected.DELETE("/tasks/:id", taskHandler.DeleteTask)

			// Context Routes
			protected.POST("/contexts", contextHandler.CreateContext)
			protected.GET("/contexts", contextHandler.ListContexts)
			protected.PATCH("/contexts/:id", contextHandler.UpdateContext)
			protected.DELETE("/contexts/:id", contextHandler.DeleteContext)
		}
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
