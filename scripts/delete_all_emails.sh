#!/usr/bin/env bash
# Script to delete all synced emails and downstream data
# Usage: ./scripts/delete_all_emails.sh

set -e

echo "⚠️  WARNING: This will delete ALL emails and related data!"
echo "This includes:"
echo "  - All emails in the 'emails' table"
echo "  - All email embeddings (automatically via CASCADE)"
echo "  - All tasks linked to emails (SourceEmailID will be set to NULL)"
echo ""
read -p "Are you sure you want to continue? (yes/no): " confirm

if [ "$confirm" != "yes" ]; then
    echo "Aborted."
    exit 0
fi

echo ""
echo "🔍 Checking database connection..."

# Get database connection details from environment or config
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_NAME=${DB_NAME:-echomind}
DB_USER=${DB_USER:-postgres}

# Check if PGPASSWORD is set
if [ -z "$PGPASSWORD" ]; then
    echo "Please set PGPASSWORD environment variable"
    echo "Example: export PGPASSWORD=your_password"
    exit 1
fi

# Test connection
if ! psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "SELECT 1" > /dev/null 2>&1; then
    echo "❌ Failed to connect to database"
    exit 1
fi

echo "✅ Database connection successful"
echo ""

# Count emails before deletion
EMAIL_COUNT=$(psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "SELECT COUNT(*) FROM emails;")
EMBEDDING_COUNT=$(psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "SELECT COUNT(*) FROM email_embeddings;")
TASK_COUNT=$(psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "SELECT COUNT(*) FROM tasks WHERE source_email_id IS NOT NULL;")

echo "📊 Current data:"
echo "  - Emails: $EMAIL_COUNT"
echo "  - Email embeddings: $EMBEDDING_COUNT"
echo "  - Tasks linked to emails: $TASK_COUNT"
echo ""

read -p "Proceed with deletion? (yes/no): " final_confirm

if [ "$final_confirm" != "yes" ]; then
    echo "Aborted."
    exit 0
fi

echo ""
echo "🗑️  Deleting data..."

# Start transaction and delete emails
# NOTE: email_embeddings will be automatically deleted via CASCADE
# Tasks with source_email_id will have that field set to NULL (no CASCADE)
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" << EOF
BEGIN;

-- Delete all emails (this will cascade to email_embeddings)
DELETE FROM emails;

-- Optional: Reset sequences if needed
-- ALTER SEQUENCE email_embeddings_id_seq RESTART WITH 1;

COMMIT;

-- Vacuum to reclaim space
VACUUM ANALYZE emails;
VACUUM ANALYZE email_embeddings;
EOF

echo ""
echo "✅ Deletion complete!"
echo ""

# Count emails after deletion
EMAIL_COUNT_AFTER=$(psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "SELECT COUNT(*) FROM emails;")
EMBEDDING_COUNT_AFTER=$(psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "SELECT COUNT(*) FROM email_embeddings;")

echo "📊 After deletion:"
echo "  - Emails: $EMAIL_COUNT_AFTER"
echo "  - Email embeddings: $EMBEDDING_COUNT_AFTER"
echo ""
echo "✨ Done!"
