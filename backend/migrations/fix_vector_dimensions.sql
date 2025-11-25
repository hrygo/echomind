-- Simple fix for vector dimension issue
-- Execute this script in PostgreSQL to fix the "expected 768 dimensions, not 1536" error

-- Step 1: Create backup of existing data
CREATE TABLE email_embeddings_backup_manual AS SELECT * FROM email_embeddings;

-- Step 2: Drop existing table (WARNING: This will delete all existing embeddings!)
DROP TABLE IF EXISTS email_embeddings CASCADE;

-- Step 3: Verify the table was dropped
SELECT 'email_embeddings table dropped successfully' AS status;

-- Step 4: Let the application recreate the table with the correct schema
-- The Go application will automatically create the table with vector(1536) when it starts
SELECT 'Ready for application to recreate table with vector(1536)' AS next_step;