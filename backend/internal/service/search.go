package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hrygo/echomind/internal/model"
	"github.com/hrygo/echomind/pkg/ai"
	"github.com/hrygo/echomind/pkg/utils"
	"github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
)

type SearchService struct {
	db       *gorm.DB
	embedder ai.EmbeddingProvider
}

func NewSearchService(db *gorm.DB, embedder ai.EmbeddingProvider) *SearchService {
	return &SearchService{
		db:       db,
		embedder: embedder,
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
	// 1. Generate embedding for the query with timeout
	embeddingCtx, cancel := context.WithTimeout(ctx, 20*time.Second) // 20 second timeout for embedding
	defer cancel()

	queryVector, err := s.embedder.Embed(embeddingCtx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to embed query: %w", err)
	}

	// 2. Validate query vector dimensions (database layer handles conversion)
	if err := s.validateVectorDimensions(queryVector, "search query"); err != nil {
		return nil, err
	}

	// 2. Perform vector search using raw SQL
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
		return nil, fmt.Errorf("search query failed: %w", err)
	}

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

	// 3. Validate vector dimensions
	if err := s.validateVectorDimensions(vectors[0], "email embedding generation"); err != nil {
		return err
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

// validateVectorDimensions validates that the vector dimensions are reasonable for processing
func (s *SearchService) validateVectorDimensions(vector []float32, context string) error {
	// Database schema supports up to 1536 dimensions with automatic conversion
	// This validation ensures vectors are reasonable for processing
	maxSupportedDimensions := 1536
	minSupportedDimensions := 1

	vectorLength := len(vector)

	if vectorLength > maxSupportedDimensions {
		return fmt.Errorf("embedding dimension too large for %s: %d dimensions (max: %d)",
			context, vectorLength, maxSupportedDimensions)
	}

	if vectorLength < minSupportedDimensions {
		return fmt.Errorf("embedding dimension too small for %s: %d dimensions (min: %d)",
			context, vectorLength, minSupportedDimensions)
	}

	// Database layer handles automatic padding/truncation in BeforeCreate hook
	// No need to validate exact dimensions here
	return nil
}
