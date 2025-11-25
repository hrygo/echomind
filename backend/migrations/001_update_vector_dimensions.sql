-- Migration: Update email_embeddings table to support dynamic vector dimensions
-- Description: Update vector column from fixed dimensions to vector(1536) to support multiple embedding providers

-- BEGIN MIGRATION

-- Step 1: Drop existing indexes on the vector column
DROP INDEX IF EXISTS email_embeddings_vector_idx;
DROP INDEX IF EXISTS email_embeddings_hnsw_idx;
DROP INDEX IF EXISTS email_embeddings_ivfflat_idx;

-- Step 2: Add new dimensions column to track actual vector dimensions
ALTER TABLE email_embeddings ADD COLUMN IF NOT EXISTS dimensions INTEGER NOT NULL DEFAULT 1024;

-- Step 3: Update existing records to set their actual dimensions based on current vector length
-- This is a best-effort approach - you may need to adjust based on your actual data

-- If most vectors are 768 dimensions (Gemini):
UPDATE email_embeddings SET dimensions = 768 WHERE array_length(vector::real[], 1) = 768;

-- If most vectors are 1024 dimensions (SiliconFlow):
UPDATE email_embeddings SET dimensions = 1024 WHERE array_length(vector::real[], 1) = 1024;

-- If most vectors are 1536 dimensions (OpenAI):
UPDATE email_embeddings SET dimensions = 1536 WHERE array_length(vector::real[], 1) = 1536;

-- Step 4: Create a backup of existing data (optional but recommended)
CREATE TABLE email_embeddings_backup AS SELECT * FROM email_embeddings;

-- Step 5: Drop the old vector column
ALTER TABLE email_embeddings DROP COLUMN IF EXISTS vector;

-- Step 6: Add the new vector column with 1536 dimensions support
ALTER TABLE email_embeddings ADD COLUMN vector vector(1536);

-- Step 7: Restore and pad/truncate vectors to 1536 dimensions
-- Pad shorter vectors with zeros, truncate longer vectors
UPDATE email_embeddings SET
  vector = CASE
    WHEN dimensions = 768 THEN vector || array_fill(0, ARRAY[768])::real[]  -- Pad 768 to 1536
    WHEN dimensions = 1024 THEN vector || array_fill(0, ARRAY[512])::real[] -- Pad 1024 to 1536
    WHEN dimensions = 1536 THEN vector                                       -- Keep 1536 as is
    WHEN dimensions > 1536 THEN vector[1:1536]                               -- Truncate if >1536
    ELSE array_fill(0, ARRAY[1536])::real[]                                  -- Default zero vector
  END;

-- Step 8: Recreate indexes with the new vector column
-- Choose the appropriate index type based on your data size and performance needs

-- For small to medium datasets (< 1M records):
CREATE INDEX email_embeddings_vector_idx
ON email_embeddings
USING ivfflat (vector vector_l2_ops)
WITH (lists = 100);

-- For large datasets (> 1M records), uncomment to use HNSW index:
-- CREATE INDEX email_embeddings_hnsw_idx
-- ON email_embeddings
-- USING hnsw (vector vector_cosine_ops)
-- WITH (m = 16, ef_construction = 64);

-- Step 9: Update table statistics
ANALYZE email_embeddings;

-- Step 10: Verify the migration
SELECT
  COUNT(*) as total_records,
  COUNT(CASE WHEN dimensions = 768 THEN 1 END) as gemini_vectors,
  COUNT(CASE WHEN dimensions = 1024 THEN 1 END) as siliconflow_vectors,
  COUNT(CASE WHEN dimensions = 1536 THEN 1 END) as openai_vectors,
  MIN(dimensions) as min_dimensions,
  MAX(dimensions) as max_dimensions
FROM email_embeddings;

-- END MIGRATION

-- Rollback instructions (if needed):
-- 1. DROP TABLE email_embeddings;
-- 2. ALTER TABLE email_embeddings_backup RENAME TO email_embeddings;
-- 3. Recreate your original indexes