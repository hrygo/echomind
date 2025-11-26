package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hrygo/echomind/internal/model"
	"github.com/hrygo/echomind/pkg/ai"
	"github.com/hrygo/echomind/pkg/telemetry"
	"github.com/hrygo/echomind/pkg/utils"
	"github.com/pgvector/pgvector-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

var tracer = otel.Tracer("echomind.search")

type SearchService struct {
	db       *gorm.DB
	embedder ai.EmbeddingProvider
	metrics  *telemetry.SearchMetrics
	cache    *SearchCache
}

func NewSearchService(db *gorm.DB, embedder ai.EmbeddingProvider, cache *SearchCache) *SearchService {
	// Initialize metrics (best effort)
	metrics, err := telemetry.NewSearchMetrics(context.Background())
	if err != nil {
		// Log error but don't fail service creation
		fmt.Printf("Warning: Failed to initialize search metrics: %v\n", err)
	}

	return &SearchService{
		db:       db,
		embedder: embedder,
		metrics:  metrics,
		cache:    cache,
	}
}

type SearchResult struct {
	EmailID uuid.UUID `json:"email_id"`
	Subject string    `json:"subject"`
	Snippet string    `json:"snippet"`
	Sender  string    `json:"sender"`
	Date    time.Time `json:"date"`
	Score   float64   `json:"score"` // Similarity score (1 - distance)
}

type SearchFilters struct {
	Sender    string
	StartDate *time.Time
	EndDate   *time.Time
	ContextID *uuid.UUID
}

func (s *SearchService) Search(ctx context.Context, userID uuid.UUID, query string, filters SearchFilters, limit int) ([]SearchResult, error) {
	// Create root span
	ctx, span := tracer.Start(ctx, "SearchService.Search",
		trace.WithAttributes(
			attribute.String("user.id", userID.String()),
			attribute.String("search.query", query),
			attribute.Int("search.limit", limit),
		),
	)
	defer span.End()

	start := time.Now()

	// Track active searches
	if s.metrics != nil {
		s.metrics.IncrementActiveSearches(ctx)
		defer s.metrics.DecrementActiveSearches(ctx)
		s.metrics.IncrementSearchRequests(ctx)
	}

	// Check cache first
	if s.cache != nil {
		cachedResults, found, err := s.cache.Get(ctx, userID, query, filters, limit)
		if err != nil {
			// Log cache error but continue with normal search
			span.AddEvent("cache_error", trace.WithAttributes(
				attribute.String("error", err.Error()),
			))
		} else if found {
			// Cache hit
			if s.metrics != nil {
				s.metrics.IncrementCacheHits(ctx)
				s.metrics.RecordSearchLatency(ctx, time.Since(start))
				s.metrics.RecordResultsReturned(ctx, len(cachedResults))
			}
			span.SetStatus(codes.Ok, "search completed (cached)")
			span.SetAttributes(
				attribute.Bool("cache.hit", true),
				attribute.Int("results.total", len(cachedResults)),
			)
			return cachedResults, nil
		}
		// Cache miss
		if s.metrics != nil {
			s.metrics.IncrementCacheMisses(ctx)
		}
		span.AddEvent("cache_miss")
	}

	// 1. Generate query embedding
	embedStart := time.Now()
	ctx, embedSpan := tracer.Start(ctx, "generate_query_embedding")
	queryVector, err := s.embedder.Embed(ctx, query)
	if err != nil {
		embedSpan.RecordError(err)
		embedSpan.SetStatus(codes.Error, "failed to embed query")
		embedSpan.End()
		span.RecordError(err)
		span.SetStatus(codes.Error, "search failed")
		if s.metrics != nil {
			s.metrics.IncrementSearchErrors(ctx)
		}
		return nil, fmt.Errorf("failed to embed query: %w", err)
	}
	embedSpan.SetAttributes(
		attribute.Int("embedding.dimensions", len(queryVector)),
	)
	embedSpan.End()
	if s.metrics != nil {
		s.metrics.RecordEmbeddingLatency(ctx, time.Since(embedStart))
	}

	// 2. Perform vector search using raw SQL
	dbStart := time.Now()
	ctx, dbSpan := tracer.Start(ctx, "vector_database_search")
	var results []SearchResult

	// Base query
	sql := `
		SELECT
			e.id as email_id,
			e.subject,
			ee.content as snippet,
			e.sender,
			e.date,
			1 - (ee.vector <=> ?) as score
		FROM email_embeddings ee
		JOIN emails e ON e.id = ee.email_id
	`
	args := []interface{}{pgvector.NewVector(queryVector)}

	// Add joins and where clauses dynamically
	var whereClauses []string
	whereClauses = append(whereClauses, "e.user_id = ?")
	args = append(args, userID)

	if filters.ContextID != nil {
		sql += " JOIN email_contexts ec ON e.id = ec.email_id"
		whereClauses = append(whereClauses, "ec.context_id = ?")
		args = append(args, filters.ContextID)
	}
	if filters.Sender != "" {
		whereClauses = append(whereClauses, "e.sender ILIKE ?")
		args = append(args, "%"+filters.Sender+"%")
	}
	if filters.StartDate != nil {
		whereClauses = append(whereClauses, "e.date >= ?")
		args = append(args, filters.StartDate)
	}
	if filters.EndDate != nil {
		whereClauses = append(whereClauses, "e.date <= ?")
		args = append(args, filters.EndDate)
	}

	if len(whereClauses) > 0 {
		sql += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Order by and Limit
	sql += " ORDER BY ee.vector <=> ? LIMIT ?"
	args = append(args, pgvector.NewVector(queryVector), limit)

	err = s.db.WithContext(ctx).Raw(sql, args...).Scan(&results).Error
	if err != nil {
		dbSpan.RecordError(err)
		dbSpan.SetStatus(codes.Error, "database query failed")
		dbSpan.End()
		span.RecordError(err)
		span.SetStatus(codes.Error, "search failed")
		if s.metrics != nil {
			s.metrics.IncrementSearchErrors(ctx)
		}
		return nil, fmt.Errorf("search query failed: %w", err)
	}

	dbSpan.SetAttributes(
		attribute.Int("results.count", len(results)),
	)
	dbSpan.End()
	if s.metrics != nil {
		s.metrics.RecordDBQueryLatency(ctx, time.Since(dbStart))
	}

	// Store in cache
	if s.cache != nil && len(results) > 0 {
		if err := s.cache.Set(ctx, userID, query, filters, limit, results); err != nil {
			// Log cache error but don't fail the request
			span.AddEvent("cache_set_error", trace.WithAttributes(
				attribute.String("error", err.Error()),
			))
		}
	}

	// Record overall metrics
	if s.metrics != nil {
		s.metrics.RecordSearchLatency(ctx, time.Since(start))
		s.metrics.RecordResultsReturned(ctx, len(results))
	}

	span.SetStatus(codes.Ok, "search completed")
	span.SetAttributes(
		attribute.Bool("cache.hit", false),
		attribute.Int("results.total", len(results)),
	)

	return results, nil
}

// GenerateAndSaveEmbedding generates and saves embeddings for an email.
func (s *SearchService) GenerateAndSaveEmbedding(ctx context.Context, email *model.Email, chunkSize int) error {
	// 1. Prepare text
	// Combine Subject and Body.
	cleanBody := utils.StripHTML(email.BodyText)
	if cleanBody == "" {
		cleanBody = email.Snippet
	}

	fullText := fmt.Sprintf("Subject: %s\n\n%s", email.Subject, cleanBody)

	// 2. Chunk text
	if chunkSize <= 0 {
		chunkSize = 1000 // Default if not specified
	}
	chunker := utils.NewTextChunker(chunkSize)
	chunks := chunker.Chunk(fullText)

	if len(chunks) == 0 {
		return nil
	}

	// 3. Generate Embeddings with timeout
	embeddingCtx, cancel := context.WithTimeout(ctx, 45*time.Second) // 45 second timeout for batch embedding
	defer cancel()

	vectors, err := s.embedder.EmbedBatch(embeddingCtx, chunks)
	if err != nil {
		return fmt.Errorf("failed to generate embeddings: %w", err)
	}

	if len(vectors) != len(chunks) {
		return fmt.Errorf("mismatch between chunks (%d) and vectors (%d)", len(chunks), len(vectors))
	}

	// 3. Basic validation - vector dimension consistency is handled by database schema
	if len(vectors) == 0 || len(vectors[0]) == 0 {
		return fmt.Errorf("no valid vectors generated")
	}

	// 4. Save to DB
	// Delete existing embeddings for this email first (to avoid duplicates on re-analysis)
	if err := s.db.WithContext(ctx).Where("email_id = ?", email.ID).Delete(&model.EmailEmbedding{}).Error; err != nil {
		return fmt.Errorf("failed to delete old embeddings: %w", err)
	}

	var embeddings []model.EmailEmbedding
	for i, vec := range vectors {
		embeddings = append(embeddings, model.EmailEmbedding{
			EmailID: email.ID,
			Content: chunks[i], // Store the actual text chunk
			Vector:  pgvector.NewVector(vec),
		})
	}

	if len(embeddings) > 0 {
		if err := s.db.WithContext(ctx).Create(&embeddings).Error; err != nil {
			return fmt.Errorf("failed to save embeddings: %w", err)
		}
	}

	return nil
}
