package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hrygo/echomind/internal/bootstrap"
	"github.com/hrygo/echomind/internal/handler"
	"github.com/hrygo/echomind/internal/middleware"
	"github.com/hrygo/echomind/internal/service"
	"github.com/hrygo/echomind/pkg/ai"
)

const Version = "0.9.0"

func main() {
	// 1. Bootstrap Application
	app, err := bootstrap.Init("configs/config.yaml", false) // false = development logger
	if err != nil {
		log.Fatalf("Bootstrap failed: %v", err)
	}
	defer app.Close()

	app.Sugar.Infof("EchoMind Version: %s", Version)

	// 2. Setup Database (Migrations & Extensions)
	if err := app.SetupDB(); err != nil {
		app.Sugar.Fatalf("Database setup failed: %v", err)
	}
	app.Sugar.Infof("Database ready")

	// 3. Initialize Services
	// Dependencies
	defaultFetcher := &service.DefaultFetcher{}

	aiProvider, err := service.NewAIProvider(&app.Config.AI)
	if err != nil {
		app.Sugar.Fatalf("Failed to create AI provider: %v", err)
	}

	organizationService := service.NewOrganizationService(app.DB)
	userService := service.NewUserService(app.DB, app.Config.Server.JWT, organizationService)
	emailService := service.NewEmailService(app.DB)
	contactService := service.NewContactService(app.DB)
	accountService := service.NewAccountService(app.DB, &app.Config.Security)
	insightService := service.NewInsightService(app.DB)
	aiDraftService := service.NewAIDraftService(aiProvider)
	syncService := service.NewSyncService(
		app.DB, 
		&service.DefaultIMAPClient{}, 
		defaultFetcher, 
		app.AsynqClient, 
		contactService, 
		accountService, 
		app.Config, 
		app.Sugar,
	)

	// Run Organization Migration
	if err := organizationService.EnsureAllUsersHaveOrganization(context.Background()); err != nil {
		app.Sugar.Errorf("Failed to migrate organizations: %v", err)
	}

	embedder, ok := aiProvider.(ai.EmbeddingProvider)
	if !ok {
		app.Sugar.Fatal("AI provider does not implement EmbeddingProvider")
	}
	searchService := service.NewSearchService(app.DB, embedder)
	chatService := service.NewChatService(aiProvider, searchService)
	taskService := service.NewTaskService(app.DB)
	contextService := service.NewContextService(app.DB)

	// 4. Initialize Handlers
	accountHandler := handler.NewAccountHandler(accountService)
	syncHandler := handler.NewSyncHandler(syncService)
	emailHandler := handler.NewEmailHandler(emailService)
	authHandler := handler.NewAuthHandler(userService)
	insightHandler := handler.NewInsightHandler(insightService)
	aiDraftHandler := handler.NewAIDraftHandler(aiDraftService)
	searchHandler := handler.NewSearchHandler(searchService, app.Sugar)
	healthHandler := handler.NewHealthHandler(app.DB)
	orgHandler := handler.NewOrganizationHandler(organizationService)
	chatHandler := handler.NewChatHandler(chatService)
	taskHandler := handler.NewTaskHandler(taskService)
	contextHandler := handler.NewContextHandler(contextService)

	// 5. Setup Router
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Register routes
	api := r.Group("/api/v1")
	{
		api.GET("/health", healthHandler.HealthCheck)
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)

		protected := api.Group("/").Use(middleware.AuthMiddleware(app.Config.Server.JWT))
		{
			// Organization
			protected.POST("/orgs", orgHandler.CreateOrganization)
			protected.GET("/orgs", orgHandler.ListOrganizations)
			protected.GET("/orgs/:id", orgHandler.GetOrganization)
			protected.GET("/orgs/:id/members", orgHandler.GetMembers)

			// Account & Sync
			protected.POST("/settings/account", accountHandler.ConnectAndSaveAccount)
			protected.GET("/settings/account", accountHandler.GetAccountStatus)
			protected.POST("/sync", syncHandler.SyncEmails)
			
			// Emails & Insights
			protected.GET("/emails", emailHandler.ListEmails)
			protected.GET("/emails/:id", emailHandler.GetEmail)
			protected.DELETE("/emails/all", emailHandler.DeleteAllEmails)
			protected.GET("/insights/network", insightHandler.GetNetworkGraph)
			
			// AI & Search
			protected.POST("/ai/draft", aiDraftHandler.GenerateDraft)
			protected.GET("/search", searchHandler.Search)
			protected.POST("/chat/completions", chatHandler.StreamChat)

			// Tasks
			protected.POST("/tasks", taskHandler.CreateTask)
			protected.GET("/tasks", taskHandler.ListTasks)
			protected.PATCH("/tasks/:id", taskHandler.UpdateTask)
			protected.PATCH("/tasks/:id/status", taskHandler.UpdateTaskStatus)
			protected.DELETE("/tasks/:id", taskHandler.DeleteTask)

			// Contexts
			protected.POST("/contexts", contextHandler.CreateContext)
			protected.GET("/contexts", contextHandler.ListContexts)
			protected.PATCH("/contexts/:id", contextHandler.UpdateContext)
			protected.DELETE("/contexts/:id", contextHandler.DeleteContext)
		}
	}

	port := app.Config.Server.Port
	if port == "" {
		port = "8080"
	}

	app.Sugar.Infof("Starting server on :%s", port)
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		app.Sugar.Fatalf("Failed to run server: %v", err)
	}
}