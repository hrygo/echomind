-- Quick fix for vector dimension mismatch
-- Use this if you're getting "expected 768 dimensions, not 1536" error

-- Step 1: Check current table structure
\d email_embeddings

-- Step 2: Add dimensions column if it doesn't exist
ALTER TABLE email_embeddings ADD COLUMN IF NOT EXISTS dimensions INTEGER NOT NULL DEFAULT 1024;

-- Step 3: Check current vector dimensions
SELECT
  id,
  array_length(vector::real[], 1) as current_dimensions
FROM email_embeddings
LIMIT 5;

-- Step 4: Update dimensions based on current vector length
UPDATE email_embeddings
SET dimensions = array_length(vector::real[], 1)
WHERE array_length(vector::real[], 1) IS NOT NULL;

-- Step 5: Create backup
CREATE TABLE email_embeddings_backup_$(date +%Y%m%d_%H%M%S) AS SELECT * FROM email_embeddings;

-- Step 6: Drop and recreate vector column with 1536 dimensions
ALTER TABLE email_embeddings DROP COLUMN vector;
ALTER TABLE email_embeddings ADD COLUMN vector vector(1536);

-- Step 7: Restore vectors with padding
UPDATE email_embeddings SET
  vector = CASE
    WHEN dimensions = 768 THEN (SELECT array_agg(elem) FROM unnest(vector_backup::real[]) WITH ORDINALITY AS t(elem, idx) UNION ALL SELECT 0 FROM generate_series(1, 768))::vector
    WHEN dimensions = 1024 THEN (SELECT array_agg(elem) FROM unnest(vector_backup::real[]) WITH ORDINALITY AS t(elem, idx) UNION ALL SELECT 0 FROM generate_series(1, 512))::vector
    WHEN dimensions = 1536 THEN vector_backup
    ELSE array_fill(0, ARRAY[1536])::real[]::vector
  END;

-- Alternative simpler approach if the above doesn't work:
-- Step 7 alternative: Just drop all data and start fresh
-- TRUNCATE TABLE email_embeddings;

-- Step 8: Verify results
SELECT
  dimensions,
  COUNT(*) as count,
  array_length(vector::real[], 1) as actual_vector_size
FROM email_embeddings
GROUP BY dimensions;

-- Step 9: Create index
CREATE INDEX IF NOT EXISTS email_embeddings_vector_idx
ON email_embeddings
USING ivfflat (vector vector_l2_ops)
WITH (lists = 100);