package app

import (
	"fmt"
	"time"

	"github.com/hrygo/echomind/internal/bootstrap"
	"github.com/hrygo/echomind/internal/event"
	"github.com/hrygo/echomind/internal/listener"
	"github.com/hrygo/echomind/internal/repository"
	"github.com/hrygo/echomind/internal/service"
	"github.com/hrygo/echomind/pkg/ai"
	"github.com/hrygo/echomind/pkg/event/bus"
	"github.com/redis/go-redis/v9"
)

// Container holds all application dependencies
// It extends the bootstrap.App with commonly used services
type Container struct {
	*bootstrap.App
	AIProvider            ai.AIProvider
	Embedder              ai.EmbeddingProvider
	SearchService         *service.SearchService
	SearchClusteringService *service.SearchClusteringService
	SearchSummaryService  *service.SearchSummaryService
	ContextService        *service.ContextService
	Summarizer            *service.SummaryService
	ActionService         *service.ActionService
	SyncService           *service.SyncService // Add SyncService
	EmailRepo             repository.EmailRepository
	AccountRepo           repository.AccountRepository
	EventBus              *bus.Bus
}

// NewContainer creates a new dependency injection container
// It initializes all common services used across cmd programs
func NewContainer(configPath string, isProduction bool) (*Container, error) {
	// 1. Bootstrap base application
	app, err := bootstrap.Init(configPath, isProduction)
	if err != nil {
		return nil, fmt.Errorf("bootstrap init failed: %w", err)
	}

	// 2. Initialize AI Provider
	aiProvider, err := service.NewAIProvider(&app.Config.AI)
	if err != nil {
		app.Close()
		return nil, fmt.Errorf("failed to create AI provider: %w", err)
	}

	// 3. Type assert Embedder
	embedder, ok := aiProvider.(ai.EmbeddingProvider)
	if !ok {
		app.Close()
		return nil, fmt.Errorf("AI provider does not implement EmbeddingProvider")
	}

	// 4. Create Redis client for caching
	var searchCache *service.SearchCache
	if app.Config.Redis.Addr != "" {
		redisClient := redis.NewClient(&redis.Options{
			Addr:     app.Config.Redis.Addr,
			Password: app.Config.Redis.Password,
			DB:       app.Config.Redis.DB,
		})
		searchCache = service.NewSearchCache(redisClient, 30*time.Minute)
	}

	// 5. Create common services
	searchService := service.NewSearchService(app.DB, embedder, searchCache)
	searchClusteringService := service.NewSearchClusteringService()
	searchSummaryService := service.NewSearchSummaryService(aiProvider)
	contextService := service.NewContextService(app.DB)
	summarizer := service.NewSummaryService(aiProvider)
	actionService := service.NewActionService(app.DB)

	// 6. Initialize Event Bus and Listeners
	eventBus := bus.New()
	contactService := service.NewContactService(app.DB) // Need this for listener
	analysisListener := listener.NewAnalysisListener(app.AsynqClient, app.Logger)
	contactListener := listener.NewContactListener(contactService, app.Logger)

	eventBus.Subscribe(event.EmailSyncedEventName, analysisListener)
	eventBus.Subscribe(event.EmailSyncedEventName, contactListener)

	// Create SyncService with dependencies
	emailRepo := repository.NewEmailRepository(app.DB)
	accountRepo := repository.NewAccountRepository(app.DB)
	imapClient := &service.DefaultIMAPClient{}

	connector := service.NewIMAPConnector(imapClient, app.Config)
	ingestor := service.NewEmailIngestor(emailRepo, app.Logger)

	syncService := service.NewSyncService(
		accountRepo,
		connector,
		ingestor,
		eventBus,
		nil, // accountService will be set later
		app.Config,
		app.Logger,
	)

	return &Container{
		App:                     app,
		AIProvider:              aiProvider,
		Embedder:                embedder,
		SearchService:           searchService,
		SearchClusteringService: searchClusteringService,
		SearchSummaryService:    searchSummaryService,
		ContextService:          contextService,
		Summarizer:              summarizer,
		ActionService:           actionService,
		SyncService:             syncService, // Add SyncService
		EmailRepo:               emailRepo,
		AccountRepo:             accountRepo,
		EventBus:                eventBus,
	}, nil
}

// ChunkSize returns the chunk size from configuration with fallback
func (c *Container) ChunkSize() int {
	if c.Config.AI.ChunkSize > 0 {
		return c.Config.AI.ChunkSize
	}
	return 1000 // Default fallback
}

// WorkerConcurrency returns the worker concurrency from configuration with fallback
func (c *Container) WorkerConcurrency() int {
	if c.Config.Worker.Concurrency > 0 {
		return c.Config.Worker.Concurrency
	}
	return 10 // Default fallback
}

// IsProduction returns true if running in production environment
func (c *Container) IsProduction() bool {
	return c.Config.Server.Environment == "production"
}
