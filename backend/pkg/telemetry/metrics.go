package telemetry

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

// SearchMetrics defines OpenTelemetry metrics for the search service
type SearchMetrics struct {
	// Histogram metrics for latency tracking
	SearchLatency    metric.Float64Histogram
	EmbeddingLatency metric.Float64Histogram
	DBQueryLatency   metric.Float64Histogram
	
	// Counter metrics for request tracking
	SearchRequests metric.Int64Counter
	SearchErrors   metric.Int64Counter
	CacheHits      metric.Int64Counter
	CacheMisses    metric.Int64Counter
	
	// UpDownCounter for active operations
	ActiveSearches metric.Int64UpDownCounter
	
	// Gauge-like counters for tracking state
	ResultsReturned metric.Int64Counter
}

// NewSearchMetrics creates a new SearchMetrics instance
func NewSearchMetrics(ctx context.Context) (*SearchMetrics, error) {
	meter := otel.Meter("echomind.search")

	// Create latency histograms
	searchLatency, err := meter.Float64Histogram(
		"search.latency",
		metric.WithDescription("Search request end-to-end latency"),
		metric.WithUnit("ms"),
	)
	if err != nil {
		return nil, err
	}

	embeddingLatency, err := meter.Float64Histogram(
		"embedding.latency",
		metric.WithDescription("AI embedding generation latency"),
		metric.WithUnit("ms"),
	)
	if err != nil {
		return nil, err
	}

	dbQueryLatency, err := meter.Float64Histogram(
		"db.query.latency",
		metric.WithDescription("Database vector search latency"),
		metric.WithUnit("ms"),
	)
	if err != nil {
		return nil, err
	}

	// Create counters
	searchRequests, err := meter.Int64Counter(
		"search.requests.total",
		metric.WithDescription("Total number of search requests"),
	)
	if err != nil {
		return nil, err
	}

	searchErrors, err := meter.Int64Counter(
		"search.errors.total",
		metric.WithDescription("Total number of search errors"),
	)
	if err != nil {
		return nil, err
	}

	cacheHits, err := meter.Int64Counter(
		"cache.hits.total",
		metric.WithDescription("Total number of cache hits"),
	)
	if err != nil {
		return nil, err
	}

	cacheMisses, err := meter.Int64Counter(
		"cache.misses.total",
		metric.WithDescription("Total number of cache misses"),
	)
	if err != nil {
		return nil, err
	}

	// Create UpDownCounter for active searches
	activeSearches, err := meter.Int64UpDownCounter(
		"search.active",
		metric.WithDescription("Current number of active search operations"),
	)
	if err != nil {
		return nil, err
	}

	resultsReturned, err := meter.Int64Counter(
		"search.results.total",
		metric.WithDescription("Total number of search results returned"),
	)
	if err != nil {
		return nil, err
	}

	return &SearchMetrics{
		SearchLatency:    searchLatency,
		EmbeddingLatency: embeddingLatency,
		DBQueryLatency:   dbQueryLatency,
		SearchRequests:   searchRequests,
		SearchErrors:     searchErrors,
		CacheHits:        cacheHits,
		CacheMisses:      cacheMisses,
		ActiveSearches:   activeSearches,
		ResultsReturned:  resultsReturned,
	}, nil
}

// RecordSearchLatency records the search latency
func (m *SearchMetrics) RecordSearchLatency(ctx context.Context, duration time.Duration) {
	m.SearchLatency.Record(ctx, float64(duration.Milliseconds()))
}

// RecordEmbeddingLatency records the embedding generation latency
func (m *SearchMetrics) RecordEmbeddingLatency(ctx context.Context, duration time.Duration) {
	m.EmbeddingLatency.Record(ctx, float64(duration.Milliseconds()))
}

// RecordDBQueryLatency records the database query latency
func (m *SearchMetrics) RecordDBQueryLatency(ctx context.Context, duration time.Duration) {
	m.DBQueryLatency.Record(ctx, float64(duration.Milliseconds()))
}

// IncrementSearchRequests increments the total search requests counter
func (m *SearchMetrics) IncrementSearchRequests(ctx context.Context) {
	m.SearchRequests.Add(ctx, 1)
}

// IncrementSearchErrors increments the search errors counter
func (m *SearchMetrics) IncrementSearchErrors(ctx context.Context) {
	m.SearchErrors.Add(ctx, 1)
}

// IncrementCacheHits increments the cache hits counter
func (m *SearchMetrics) IncrementCacheHits(ctx context.Context) {
	m.CacheHits.Add(ctx, 1)
}

// IncrementCacheMisses increments the cache misses counter
func (m *SearchMetrics) IncrementCacheMisses(ctx context.Context) {
	m.CacheMisses.Add(ctx, 1)
}

// IncrementActiveSearches increments the active searches counter
func (m *SearchMetrics) IncrementActiveSearches(ctx context.Context) {
	m.ActiveSearches.Add(ctx, 1)
}

// DecrementActiveSearches decrements the active searches counter
func (m *SearchMetrics) DecrementActiveSearches(ctx context.Context) {
	m.ActiveSearches.Add(ctx, -1)
}

// RecordResultsReturned records the number of results returned
func (m *SearchMetrics) RecordResultsReturned(ctx context.Context, count int) {
	m.ResultsReturned.Add(ctx, int64(count))
}

// SyncMetrics defines OpenTelemetry metrics for the email sync service
type SyncMetrics struct {
	SyncLatency     metric.Float64Histogram
	SyncRequests    metric.Int64Counter
	SyncErrors      metric.Int64Counter
	EmailsProcessed metric.Int64Counter
	EmailsFailed    metric.Int64Counter
}

// NewSyncMetrics creates a new SyncMetrics instance
func NewSyncMetrics(ctx context.Context) (*SyncMetrics, error) {
	meter := otel.Meter("echomind.sync")

	syncLatency, err := meter.Float64Histogram(
		"sync.latency",
		metric.WithDescription("Email sync operation latency"),
		metric.WithUnit("ms"),
	)
	if err != nil {
		return nil, err
	}

	syncRequests, err := meter.Int64Counter(
		"sync.requests.total",
		metric.WithDescription("Total number of sync requests"),
	)
	if err != nil {
		return nil, err
	}

	syncErrors, err := meter.Int64Counter(
		"sync.errors.total",
		metric.WithDescription("Total number of sync errors"),
	)
	if err != nil {
		return nil, err
	}

	emailsProcessed, err := meter.Int64Counter(
		"sync.emails.processed",
		metric.WithDescription("Total number of emails processed"),
	)
	if err != nil {
		return nil, err
	}

	emailsFailed, err := meter.Int64Counter(
		"sync.emails.failed",
		metric.WithDescription("Total number of emails failed to process"),
	)
	if err != nil {
		return nil, err
	}

	return &SyncMetrics{
		SyncLatency:     syncLatency,
		SyncRequests:    syncRequests,
		SyncErrors:      syncErrors,
		EmailsProcessed: emailsProcessed,
		EmailsFailed:    emailsFailed,
	}, nil
}

// AIMetrics defines OpenTelemetry metrics for AI services
type AIMetrics struct {
	AIRequestLatency metric.Float64Histogram
	AIRequests       metric.Int64Counter
	AIErrors         metric.Int64Counter
	TokensUsed       metric.Int64Counter
}

// NewAIMetrics creates a new AIMetrics instance
func NewAIMetrics(ctx context.Context) (*AIMetrics, error) {
	meter := otel.Meter("echomind.ai")

	aiRequestLatency, err := meter.Float64Histogram(
		"ai.request.latency",
		metric.WithDescription("AI API request latency"),
		metric.WithUnit("ms"),
	)
	if err != nil {
		return nil, err
	}

	aiRequests, err := meter.Int64Counter(
		"ai.requests.total",
		metric.WithDescription("Total number of AI API requests"),
	)
	if err != nil {
		return nil, err
	}

	aiErrors, err := meter.Int64Counter(
		"ai.errors.total",
		metric.WithDescription("Total number of AI API errors"),
	)
	if err != nil {
		return nil, err
	}

	tokensUsed, err := meter.Int64Counter(
		"ai.tokens.used",
		metric.WithDescription("Total number of AI tokens consumed"),
	)
	if err != nil {
		return nil, err
	}

	return &AIMetrics{
		AIRequestLatency: aiRequestLatency,
		AIRequests:       aiRequests,
		AIErrors:         aiErrors,
		TokensUsed:       tokensUsed,
	}, nil
}

// CacheMetrics defines OpenTelemetry metrics for cache operations
type CacheMetrics struct {
	// Latency histograms
	GetLatency    metric.Float64Histogram
	SetLatency    metric.Float64Histogram
	DeleteLatency metric.Float64Histogram
	
	// Operation counters
	Operations metric.Int64Counter
	Errors     metric.Int64Counter
	
	// Size histograms
	KeySize   metric.Int64Histogram
	ValueSize metric.Int64Histogram
	
	// Hit/Miss counters (specific to cache)
	Hits   metric.Int64Counter
	Misses metric.Int64Counter
}

// NewCacheMetrics creates a new CacheMetrics instance
func NewCacheMetrics(ctx context.Context) (*CacheMetrics, error) {
	meter := otel.Meter("echomind.cache")

	// Create latency histograms
	getLatency, err := meter.Float64Histogram(
		"cache.get.latency",
		metric.WithDescription("Cache get operation latency"),
		metric.WithUnit("ms"),
	)
	if err != nil {
		return nil, err
	}

	setLatency, err := meter.Float64Histogram(
		"cache.set.latency",
		metric.WithDescription("Cache set operation latency"),
		metric.WithUnit("ms"),
	)
	if err != nil {
		return nil, err
	}

	deleteLatency, err := meter.Float64Histogram(
		"cache.delete.latency",
		metric.WithDescription("Cache delete operation latency"),
		metric.WithUnit("ms"),
	)
	if err != nil {
		return nil, err
	}

	// Create operation counters
	operations, err := meter.Int64Counter(
		"cache.operations.total",
		metric.WithDescription("Total number of cache operations"),
	)
	if err != nil {
		return nil, err
	}

	errors, err := meter.Int64Counter(
		"cache.errors.total",
		metric.WithDescription("Total number of cache errors"),
	)
	if err != nil {
		return nil, err
	}

	// Create size histograms
	keySize, err := meter.Int64Histogram(
		"cache.key.size",
		metric.WithDescription("Cache key size distribution"),
		metric.WithUnit("bytes"),
	)
	if err != nil {
		return nil, err
	}

	valueSize, err := meter.Int64Histogram(
		"cache.value.size",
		metric.WithDescription("Cache value size distribution"),
		metric.WithUnit("bytes"),
	)
	if err != nil {
		return nil, err
	}

	// Create hit/miss counters
	hits, err := meter.Int64Counter(
		"cache.hits.total",
		metric.WithDescription("Total number of cache hits"),
	)
	if err != nil {
		return nil, err
	}

	misses, err := meter.Int64Counter(
		"cache.misses.total",
		metric.WithDescription("Total number of cache misses"),
	)
	if err != nil {
		return nil, err
	}

	return &CacheMetrics{
		GetLatency:    getLatency,
		SetLatency:    setLatency,
		DeleteLatency: deleteLatency,
		Operations:    operations,
		Errors:        errors,
		KeySize:       keySize,
		ValueSize:     valueSize,
		Hits:          hits,
		Misses:        misses,
	}, nil
}

// RecordGetLatency records cache get operation latency
func (m *CacheMetrics) RecordGetLatency(ctx context.Context, latencyMs int64) {
	if m.GetLatency != nil {
		m.GetLatency.Record(ctx, float64(latencyMs))
	}
}

// RecordSetLatency records cache set operation latency
func (m *CacheMetrics) RecordSetLatency(ctx context.Context, latencyMs int64) {
	if m.SetLatency != nil {
		m.SetLatency.Record(ctx, float64(latencyMs))
	}
}

// RecordDeleteLatency records cache delete operation latency
func (m *CacheMetrics) RecordDeleteLatency(ctx context.Context, latencyMs int64) {
	if m.DeleteLatency != nil {
		m.DeleteLatency.Record(ctx, float64(latencyMs))
	}
}

// IncrementOperations increments the operations counter
func (m *CacheMetrics) IncrementOperations(ctx context.Context, operation string) {
	if m.Operations != nil {
		m.Operations.Add(ctx, 1)
	}
}

// IncrementErrors increments the errors counter
func (m *CacheMetrics) IncrementErrors(ctx context.Context, errorType string) {
	if m.Errors != nil {
		m.Errors.Add(ctx, 1)
	}
}

// RecordKeySize records cache key size
func (m *CacheMetrics) RecordKeySize(ctx context.Context, size int64) {
	if m.KeySize != nil {
		m.KeySize.Record(ctx, size)
	}
}

// RecordValueSize records cache value size
func (m *CacheMetrics) RecordValueSize(ctx context.Context, size int64) {
	if m.ValueSize != nil {
		m.ValueSize.Record(ctx, size)
	}
}

// IncrementHits increments the cache hits counter
func (m *CacheMetrics) IncrementHits(ctx context.Context) {
	if m.Hits != nil {
		m.Hits.Add(ctx, 1)
	}
}

// IncrementMisses increments the cache misses counter
func (m *CacheMetrics) IncrementMisses(ctx context.Context) {
	if m.Misses != nil {
		m.Misses.Add(ctx, 1)
	}
}
