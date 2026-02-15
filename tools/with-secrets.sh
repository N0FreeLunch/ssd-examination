#!/bin/bash
# Location: tools/with-secrets.sh
# Purpose: Fetch secrets from AWS SSM and execute command with injected environment variables.
# Usage: ./tools/with-secrets.sh [local|prod] [command]

ENV=$1
shift
CMD="$@"

if [ -z "$ENV" ] || [ -z "$CMD" ]; then
    echo "❌ Usage: $0 [local|prod] [command]"
    exit 1
fi

echo "🔐 Loading secrets for '$ENV' from AWS Parameter Store..."

# Parameter path prefix based on environment
PREFIX="/sdd-exam/$ENV"

# Fetch parameters from AWS SSM (Recursive with Decryption)
PARAMETERS=$(aws ssm get-parameters-by-path \
    --path "$PREFIX" \
    --recursive \
    --with-decryption \
    --query "Parameters[*].{Name:Name,Value:Value}" \
    --output json)

if [ $? -ne 0 ]; then
    echo "❌ Failed to fetch parameters from AWS SSM. Please check your AWS credentials (aws configure)."
    exit 1
fi

# Check if parameters are found
if [ -z "$PARAMETERS" ] || [ "$PARAMETERS" == "[]" ]; then
    echo "⚠️  Warning: No parameters found in path: $PREFIX"
else
    # Parse JSON and export variables using jq
    # Extracts the last part of the name (e.g., /sdd-exam/local/DATABASE_URL -> DATABASE_URL)
    while IFS="=" read -r key value; do
        export "$key"="$value"
        echo "   + $key set"
    done < <(echo "$PARAMETERS" | jq -r '.[] | "\(.Name | split("/") | last)=\(.Value)"')
fi

echo "🚀 Executing command: $CMD"
exec env $CMD
