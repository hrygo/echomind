package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hrygo/echomind/internal/app"
	"github.com/hrygo/echomind/internal/handler"
	"github.com/hrygo/echomind/internal/router"
	"github.com/hrygo/echomind/internal/service"
)

const Version = "0.9.2"

func main() {
	// Parse CLI configuration
	cli := app.ParseCLI()

	// Initialize application container
	container, err := app.NewContainer(cli.ConfigPath, cli.IsProduction)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	defer container.Close()

	container.Sugar.Infof("EchoMind Version: %s", Version)

	// Setup Database (Migrations & Extensions)
	if err := container.SetupDB(); err != nil {
		container.Sugar.Fatalf("Database setup failed: %v", err)
	}
	container.Sugar.Infof("Database ready")

	// Log AI Provider Configuration
	container.Sugar.Infof("AI Provider Initialized:")
	container.Sugar.Infof("  Chat Provider: %s", container.Config.AI.ActiveServices.Chat)
	if chatProviderConfig, ok := container.Config.AI.Providers[container.Config.AI.ActiveServices.Chat]; ok {
		if model, exists := chatProviderConfig.Settings["model"]; exists {
			container.Sugar.Infof("    Model: %s", model)
		}
		if baseURL, exists := chatProviderConfig.Settings["base_url"]; exists {
			container.Sugar.Infof("    Base URL: %s", baseURL)
		}
	}
	container.Sugar.Infof("  Embedding Provider: %s", container.Config.AI.ActiveServices.Embedding)
	if embedProviderConfig, ok := container.Config.AI.Providers[container.Config.AI.ActiveServices.Embedding]; ok {
		if embeddingModel, exists := embedProviderConfig.Settings["embedding_model"]; exists {
			container.Sugar.Infof("    Embedding Model: %s", embeddingModel)
		}
	}

	// Initialize business services
	defaultFetcher := &service.DefaultFetcher{}
	organizationService := service.NewOrganizationService(container.DB)
	userService := service.NewUserService(container.DB, container.Config.Server.JWT, organizationService)
	emailService := service.NewEmailService(container.DB)
	contactService := service.NewContactService(container.DB)
	accountService := service.NewAccountService(container.DB, &container.Config.Security)
	insightService := service.NewInsightService(container.DB)
	aiDraftService := service.NewAIDraftService(container.AIProvider)
	syncService := service.NewSyncService(
		container.DB,
		&service.DefaultIMAPClient{},
		defaultFetcher,
		container.AsynqClient,
		contactService,
		accountService,
		container.Config,
		container.Sugar,
	)

	// Run Organization Migration
	if err := organizationService.EnsureAllUsersHaveOrganization(context.Background()); err != nil {
		container.Sugar.Errorf("Failed to migrate organizations: %v", err)
	}

		chatService := service.NewChatService(container.AIProvider, container.SearchService)
		taskService := service.NewTaskService(container.DB)

		// Initialize handlers
		accountHandler := handler.NewAccountHandler(accountService)
		syncHandler := handler.NewSyncHandler(syncService)
		emailHandler := handler.NewEmailHandler(emailService)
		authHandler := handler.NewAuthHandler(userService)
		insightHandler := handler.NewInsightHandler(insightService)
		aiDraftHandler := handler.NewAIDraftHandler(aiDraftService)
		searchHandler := handler.NewSearchHandler(container.SearchService, container.Sugar)
		healthHandler := handler.NewHealthHandler(container.DB)
		orgHandler := handler.NewOrganizationHandler(organizationService)
		chatHandler := handler.NewChatHandler(chatService)
		taskHandler := handler.NewTaskHandler(taskService)
		contextHandler := handler.NewContextHandler(container.ContextService)
		actionHandler := handler.NewActionHandler(container.ActionService)

		// Setup Router and Middleware
		r := gin.Default()
		router.SetupMiddleware(r, container.IsProduction())

		// Register routes
		handlers := &router.Handlers{
			Health:  healthHandler,
			Auth:    authHandler,
			Org:     orgHandler,
			Account: accountHandler,
			Sync:    syncHandler,
			Email:   emailHandler,
			Insight: insightHandler,
			AIDraft: aiDraftHandler,
			Search:  searchHandler,
			Chat:    chatHandler,
			Task:    taskHandler,
			Context: contextHandler,
			Action:  actionHandler,
		}

		authMiddleware := router.SetupAuthMiddleware(container.Config.Server.JWT)
		router.SetupRoutes(r, handlers, authMiddleware)

		port := container.Config.Server.Port

	// Create HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		container.Sugar.Infof("Starting server on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			container.Sugar.Fatalf("Failed to run server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	container.Sugar.Info("Shutting down server...")

	// Give outstanding requests 10 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		container.Sugar.Fatalf("Server forced to shutdown: %v", err)
	}

	container.Sugar.Info("Server exited gracefully")
}
