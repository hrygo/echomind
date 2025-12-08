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

const Version = "0.9.8"

func main() {
	// Parse CLI configuration
	cli := app.ParseCLI()

	// Initialize application container
	container, err := app.NewContainer(cli.ConfigPath, cli.IsProduction)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	defer container.Close()

	// Setup Database (Migrations & Extensions)
	if err := container.SetupDB(); err != nil {
		log.Printf("Failed to setup DB: %v", err)
	}

	// Log AI Provider Configuration
	_ = container.Config.AI.Providers[container.Config.AI.ActiveServices.Chat]

	// Initialize business services
	organizationService := service.NewOrganizationService(container.DB)
	userService := service.NewUserService(container.DB, container.Config.Server.JWT, organizationService)
	emailService := service.NewEmailService(container.DB)
	accountService := service.NewAccountService(container.DB, &container.Config.Security)
	insightService := service.NewInsightService(container.DB)
	aiDraftService := service.NewAIDraftService(container.AIProvider)
	defaultIMAPClient := &service.DefaultIMAPClient{}
	connector := service.NewIMAPConnector(defaultIMAPClient, container.Config)
	ingestor := service.NewEmailIngestor(container.EmailRepo, container.Logger)

	syncService := service.NewSyncService(
		container.AccountRepo,
		connector,
		ingestor,
		container.EventBus,
		accountService,
		container.Config,
		container.Logger,
	)

	// Run Organization Migration
	if err := organizationService.EnsureAllUsersHaveOrganization(context.Background()); err != nil {
		log.Printf("Failed to ensure organization: %v", err)
	}

	chatService := service.NewChatService(container.AIProvider, container.SearchService, emailService)
	taskService := service.NewTaskService(container.DB)
	opportunityService := service.NewOpportunityService(container.DB)

	// Initialize handlers
	accountHandler := handler.NewAccountHandler(accountService)
	syncHandler := handler.NewSyncHandler(syncService)
	emailHandler := handler.NewEmailHandler(emailService)
	authHandler := handler.NewAuthHandler(userService)
	insightHandler := handler.NewInsightHandlerWithServices(insightService, taskService, emailService)
	aiDraftHandler := handler.NewAIDraftHandler(aiDraftService, emailService)
	searchHandler := handler.NewSearchHandler(container.SearchService, container.SearchClusteringService, container.SearchSummaryService, container.Logger)
	healthHandler := handler.NewHealthHandler(container.DB)
	orgHandler := handler.NewOrganizationHandler(organizationService)
	chatHandler := handler.NewChatHandler(chatService)
	taskHandler := handler.NewTaskHandler(taskService)
	contextHandler := handler.NewContextHandler(container.ContextService)
	actionHandler := handler.NewActionHandler(container.ActionService)
	opportunityHandler := handler.NewOpportunityHandler(opportunityService)

	// Setup Router and Middleware
	r := gin.Default()
	router.SetupMiddleware(r, container.App, container.IsProduction())

	// Register routes
	handlers := &router.Handlers{
		Health:      healthHandler,
		Auth:        authHandler,
		Org:         orgHandler,
		Account:     accountHandler,
		Sync:        syncHandler,
		Email:       emailHandler,
		Insight:     insightHandler,
		AIDraft:     aiDraftHandler,
		Search:      searchHandler,
		Chat:        chatHandler,
		Task:        taskHandler,
		Context:     contextHandler,
		Action:      actionHandler,
		Opportunity: opportunityHandler,
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
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Give outstanding requests 10 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown failed: %v", err)
	}

}
