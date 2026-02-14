#!/bin/sh
set -e

DB_PATH="/data/local.db"

if [ -f "$DB_PATH" ]; then
    echo "✅ Database found at $DB_PATH. Skipping restore."
else
    echo "⚠️ Database not found at $DB_PATH."
    echo "🔄 Attempting to restore from MinIO..."
    
    litestream databases -config /etc/litestream.yml
    
    # Removed -if-replica-exists (for debugging purposes)
    litestream restore -config /etc/litestream.yml "$DB_PATH"
    
    echo "✅ Restore process completed."
fi
