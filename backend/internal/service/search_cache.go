package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hrygo/echomind/pkg/telemetry"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// SearchCache provides caching for search results with OpenTelemetry instrumentation
type SearchCache struct {
	redis   *redis.Client
	ttl     time.Duration
	metrics *telemetry.CacheMetrics
	tracer  trace.Tracer
}

// NewSearchCache creates a new search cache instance with OTel instrumentation
func NewSearchCache(redisClient *redis.Client, ttl time.Duration) *SearchCache {
	if ttl == 0 {
		ttl = 30 * time.Minute // Default 30 minutes
	}

	// Initialize metrics (best effort)
	metrics, err := telemetry.NewCacheMetrics(context.Background())
	if err != nil {
		fmt.Printf("Warning: Failed to initialize cache metrics: %v\n", err)
	}

	return &SearchCache{
		redis:   redisClient,
		ttl:     ttl,
		metrics: metrics,
		tracer:  otel.Tracer("echomind.cache"),
	}
}

// generateCacheKey generates a unique cache key for search parameters with tracing
func (c *SearchCache) generateCacheKey(ctx context.Context, userID uuid.UUID, query string, filters SearchFilters, limit int) string {
	ctx, span := c.tracer.Start(ctx, "generate_cache_key",
		trace.WithSpanKind(trace.SpanKindInternal),
	)
	defer span.End()

	// Create a deterministic string from search parameters
	filterStr := fmt.Sprintf("%s|%v|%v|%v",
		filters.Sender,
		filters.StartDate,
		filters.EndDate,
		filters.ContextID,
	)

	keyData := fmt.Sprintf("search:%s:%s:%s:%d", userID.String(), query, filterStr, limit)

	// Hash the key to keep it short
	hash := sha256.Sum256([]byte(keyData))
	key := "search:cache:" + hex.EncodeToString(hash[:16])

	// Record span attributes
	span.SetAttributes(
		attribute.String("cache.key", key),
		attribute.Int("cache.key_size", len(key)),
		attribute.String("user.id", userID.String()),
		attribute.String("search.query", query),
	)

	// Record key size metric
	if c.metrics != nil {
		c.metrics.RecordKeySize(ctx, int64(len(key)))
	}

	return key
}

// Get retrieves cached search results with full OTel instrumentation
func (c *SearchCache) Get(ctx context.Context, userID uuid.UUID, query string, filters SearchFilters, limit int) ([]SearchResult, bool, error) {
	// Create span for cache get operation
	ctx, span := c.tracer.Start(ctx, "SearchCache.Get",
		trace.WithSpanKind(trace.SpanKindInternal),
	)
	defer span.End()

	start := time.Now()

	if c.redis == nil {
		return nil, false, nil
	}

	// Generate cache key with sub-span
	key := c.generateCacheKey(ctx, userID, query, filters, limit)
	span.SetAttributes(
		attribute.String("cache.key", key),
		attribute.String("cache.operation", "get"),
		attribute.String("cache.service", "search_cache"),
		attribute.String("cache.backend", "redis"),
	)

	// Redis operation with sub-span
	ctx2, redisSpan := c.tracer.Start(ctx, "redis_get")
	data, err := c.redis.Get(ctx2, key).Bytes()
	redisSpan.End()

	latency := time.Since(start)

	// Record metrics
	if c.metrics != nil {
		c.metrics.RecordGetLatency(ctx, latency.Milliseconds())
		c.metrics.IncrementOperations(ctx, "get")
	}

	if err == redis.Nil {
		// Cache miss
		span.SetAttributes(attribute.Bool("cache.hit", false))
		if c.metrics != nil {
			c.metrics.IncrementMisses(ctx)
		}
		return nil, false, nil
	}

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "redis get error")
		if c.metrics != nil {
			c.metrics.IncrementErrors(ctx, "get")
		}
		return nil, false, fmt.Errorf("redis get error: %w", err)
	}

	// Cache hit - unmarshal results
	var results []SearchResult
	if err := json.Unmarshal(data, &results); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "unmarshal error")
		if c.metrics != nil {
			c.metrics.IncrementErrors(ctx, "unmarshal")
		}
		return nil, false, fmt.Errorf("failed to unmarshal cached results: %w", err)
	}

	// Success - record attributes and metrics
	span.SetAttributes(
		attribute.Bool("cache.hit", true),
		attribute.Int("cache.result_count", len(results)),
		attribute.Int("cache.value_size", len(data)),
	)

	if c.metrics != nil {
		c.metrics.IncrementHits(ctx)
		c.metrics.RecordValueSize(ctx, int64(len(data)))
	}

	return results, true, nil
}

// Set stores search results in cache with full OTel instrumentation
func (c *SearchCache) Set(ctx context.Context, userID uuid.UUID, query string, filters SearchFilters, limit int, results []SearchResult) error {
	// Create span for cache set operation
	ctx, span := c.tracer.Start(ctx, "SearchCache.Set",
		trace.WithSpanKind(trace.SpanKindInternal),
	)
	defer span.End()

	start := time.Now()

	if c.redis == nil {
		return nil
	}

	// Generate cache key with sub-span
	key := c.generateCacheKey(ctx, userID, query, filters, limit)
	span.SetAttributes(
		attribute.String("cache.key", key),
		attribute.String("cache.operation", "set"),
		attribute.String("cache.service", "search_cache"),
		attribute.String("cache.backend", "redis"),
		attribute.Int64("cache.ttl_seconds", int64(c.ttl.Seconds())),
	)

	// Marshal results
	data, err := json.Marshal(results)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "marshal error")
		if c.metrics != nil {
			c.metrics.IncrementErrors(ctx, "marshal")
		}
		return fmt.Errorf("failed to marshal results: %w", err)
	}

	span.SetAttributes(
		attribute.Int("cache.value_size", len(data)),
		attribute.Int("cache.result_count", len(results)),
	)

	// Redis operation with sub-span
	ctx2, redisSpan := c.tracer.Start(ctx, "redis_set")
	err = c.redis.Set(ctx2, key, data, c.ttl).Err()
	redisSpan.End()

	latency := time.Since(start)

	// Record metrics
	if c.metrics != nil {
		c.metrics.RecordSetLatency(ctx, latency.Milliseconds())
		c.metrics.IncrementOperations(ctx, "set")
		c.metrics.RecordValueSize(ctx, int64(len(data)))
	}

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "redis set error")
		if c.metrics != nil {
			c.metrics.IncrementErrors(ctx, "set")
		}
		return fmt.Errorf("redis set error: %w", err)
	}

	return nil
}

// Invalidate removes cached results for a user with full OTel instrumentation
func (c *SearchCache) Invalidate(ctx context.Context, userID uuid.UUID) error {
	// Create span for invalidate operation
	ctx, span := c.tracer.Start(ctx, "SearchCache.Invalidate",
		trace.WithSpanKind(trace.SpanKindInternal),
	)
	defer span.End()

	start := time.Now()

	if c.redis == nil {
		return nil
	}

	span.SetAttributes(
		attribute.String("cache.operation", "invalidate"),
		attribute.String("cache.service", "search_cache"),
		attribute.String("cache.backend", "redis"),
		attribute.String("user.id", userID.String()),
	)

	// Delete all search cache entries for this user
	pattern := fmt.Sprintf("search:cache:*%s*", userID.String())

	var deletedCount int
	var cursor uint64
	for {
		var keys []string
		var err error

		// Redis scan with sub-span
		ctx2, scanSpan := c.tracer.Start(ctx, "redis_scan")
		keys, cursor, err = c.redis.Scan(ctx2, cursor, pattern, 100).Result()
		scanSpan.End()

		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "redis scan error")
			if c.metrics != nil {
				c.metrics.IncrementErrors(ctx, "scan")
			}
			return fmt.Errorf("redis scan error: %w", err)
		}

		if len(keys) > 0 {
			// Redis delete with sub-span
			ctx3, delSpan := c.tracer.Start(ctx, "redis_del")
			deleted, err := c.redis.Del(ctx3, keys...).Result()
			delSpan.End()

			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, "redis del error")
				if c.metrics != nil {
					c.metrics.IncrementErrors(ctx, "del")
				}
				return fmt.Errorf("redis del error: %w", err)
			}
			deletedCount += int(deleted)
		}

		if cursor == 0 {
			break
		}
	}

	latency := time.Since(start)

	// Record attributes and metrics
	span.SetAttributes(attribute.Int("cache.keys_deleted", deletedCount))

	if c.metrics != nil {
		c.metrics.RecordDeleteLatency(ctx, latency.Milliseconds())
		c.metrics.IncrementOperations(ctx, "invalidate")
	}

	return nil
}

// InvalidateAll clears all search cache with full OTel instrumentation
func (c *SearchCache) InvalidateAll(ctx context.Context) error {
	// Create span for invalidate all operation
	ctx, span := c.tracer.Start(ctx, "SearchCache.InvalidateAll",
		trace.WithSpanKind(trace.SpanKindInternal),
	)
	defer span.End()

	start := time.Now()

	if c.redis == nil {
		return nil
	}

	span.SetAttributes(
		attribute.String("cache.operation", "invalidate_all"),
		attribute.String("cache.service", "search_cache"),
		attribute.String("cache.backend", "redis"),
	)

	pattern := "search:cache:*"

	var deletedCount int
	var cursor uint64
	for {
		var keys []string
		var err error

		// Redis scan with sub-span
		ctx2, scanSpan := c.tracer.Start(ctx, "redis_scan")
		keys, cursor, err = c.redis.Scan(ctx2, cursor, pattern, 100).Result()
		scanSpan.End()

		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "redis scan error")
			if c.metrics != nil {
				c.metrics.IncrementErrors(ctx, "scan")
			}
			return fmt.Errorf("redis scan error: %w", err)
		}

		if len(keys) > 0 {
			// Redis delete with sub-span
			ctx3, delSpan := c.tracer.Start(ctx, "redis_del")
			deleted, err := c.redis.Del(ctx3, keys...).Result()
			delSpan.End()

			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, "redis del error")
				if c.metrics != nil {
					c.metrics.IncrementErrors(ctx, "del")
				}
				return fmt.Errorf("redis del error: %w", err)
			}
			deletedCount += int(deleted)
		}

		if cursor == 0 {
			break
		}
	}

	latency := time.Since(start)

	// Record attributes and metrics
	span.SetAttributes(attribute.Int("cache.keys_deleted", deletedCount))

	if c.metrics != nil {
		c.metrics.RecordDeleteLatency(ctx, latency.Milliseconds())
		c.metrics.IncrementOperations(ctx, "invalidate_all")
	}

	return nil
}
