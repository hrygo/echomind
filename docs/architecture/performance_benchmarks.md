# Performance Benchmarks & Optimization

## Baseline Metrics (v0.6.5)

**Date**: 2025-11-22  
**Hardware**: Local Development Environment (Apple M1)  
**Database**: Postgres 16 + pgvector  

### Search Latency (Benchmarks)

| Dataset Size | Avg Latency (p95) | Ops/sec | Notes |
|--------------|-------------------|---------|-------|
| 100 emails   | ~1.67 ms          | ~600    | In-memory/cached range |
| 1,000 emails | ~9.30 ms          | ~107    | Linear scaling observed |

*Note: Extrapolated for 10k emails ~ 93ms, which is well within the < 500ms target.*

### Optimization Actions Taken

1.  **Vector Indexing**:
    - HNSW Index (`email_embeddings_vector_idx`) is applied on `email_embeddings.vector`.
    - Operator: `vector_cosine_ops` (Cosine Distance).
    - Verified via `backend/cmd/main.go`.

2.  **Chunking Strategy**:
    - Implemented configurable `chunk_size` in `AIConfig`.
    - Default: `1000` tokens (~4000 chars).
    - Recommendation: For finer granularity, lower to `500` tokens in `configs/config.yaml`.

3.  **Database Tuning**:
    - Confirmed usage of `pgvector` HNSW index.
    - Future tuning: Adjust `m` and `ef_construction` if recall drops or latency increases > 500ms.

## Next Steps

- Monitor real-world performance with larger datasets (10k+).
- Consider "Hybrid Search" (Keyword + Vector) if semantic search misses exact matches.
