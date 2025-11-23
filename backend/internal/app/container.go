package app

import (
	"fmt"

	"github.com/hrygo/echomind/internal/bootstrap"
	"github.com/hrygo/echomind/internal/service"
	"github.com/hrygo/echomind/pkg/ai"
)

// Container holds all application dependencies
// It extends the bootstrap.App with commonly used services
type Container struct {
	*bootstrap.App
	AIProvider     ai.AIProvider
	Embedder       ai.EmbeddingProvider
	SearchService  *service.SearchService
	ContextService *service.ContextService
	Summarizer     *service.SummaryService
	ActionService  *service.ActionService
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

	// 4. Create common services
	searchService := service.NewSearchService(app.DB, embedder)
	contextService := service.NewContextService(app.DB)
	summarizer := service.NewSummaryService(aiProvider)
	actionService := service.NewActionService(app.DB)

	return &Container{
		App:            app,
		AIProvider:     aiProvider,
		Embedder:       embedder,
		SearchService:  searchService,
		ContextService: contextService,
		Summarizer:     summarizer,
		ActionService:  actionService,
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
