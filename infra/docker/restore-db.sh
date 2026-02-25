#!/bin/sh
set -e

DB_PATH="/data/local.db"

if [ -f "$DB_PATH" ]; then
    echo "✅ Database found at $DB_PATH. Skipping restore."
else
    echo "⚠️ Database not found at $DB_PATH."
    echo "🔄 Attempting to restore from S3 (${LITESTREAM_BUCKET})..."
    
    # -if-replica-exists: proceed without error if backup doesn't exist (handle initial run)
    litestream restore -config /etc/litestream.yml -if-replica-exists "$DB_PATH"
    
    echo "✅ Restore process completed."
fi

