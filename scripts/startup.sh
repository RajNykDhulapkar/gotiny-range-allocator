#!/bin/sh
# scripts/startup.sh
set -e

echo "Starting Range Allocator Service..."

# Function to check required environment variables
check_required_vars() {
  local missing_vars=0

  # Database URL
  if [ -z "${RANGE_ALLOCATOR_DATABASE_URL}" ]; then
    echo "ERROR: Required environment variable RANGE_ALLOCATOR_DATABASE_URL is not set"
    missing_vars=1
  fi

  # GRPC Port
  if [ -z "${RANGE_ALLOCATOR_GRPC_PORT}" ]; then
    echo "ERROR: Required environment variable RANGE_ALLOCATOR_GRPC_PORT is not set"
    missing_vars=1
  else
    # Validate port number
    if ! echo "${RANGE_ALLOCATOR_GRPC_PORT}" | grep -Eq '^[0-9]+$'; then
      echo "ERROR: RANGE_ALLOCATOR_GRPC_PORT must be a valid port number"
      missing_vars=1
    elif [ "${RANGE_ALLOCATOR_GRPC_PORT}" -lt 1 ] || [ "${RANGE_ALLOCATOR_GRPC_PORT}" -gt 65535 ]; then
      echo "ERROR: RANGE_ALLOCATOR_GRPC_PORT must be between 1 and 65535"
      missing_vars=1
    fi
  fi

  # Range Configuration
  if [ -z "${RANGE_ALLOCATOR_RANGE_DEFAULT_SIZE}" ]; then
    echo "ERROR: Required environment variable RANGE_ALLOCATOR_RANGE_DEFAULT_SIZE is not set"
    missing_vars=1
  fi

  if [ -z "${RANGE_ALLOCATOR_RANGE_MIN_SIZE}" ]; then
    echo "ERROR: Required environment variable RANGE_ALLOCATOR_RANGE_MIN_SIZE is not set"
    missing_vars=1
  fi

  if [ -z "${RANGE_ALLOCATOR_RANGE_MAX_SIZE}" ]; then
    echo "ERROR: Required environment variable RANGE_ALLOCATOR_RANGE_MAX_SIZE is not set"
    missing_vars=1
  fi

  # Exit if any required variables are missing
  if [ $missing_vars -ne 0 ]; then
    echo "Missing required environment variables. Please set them and try again."
    exit 1
  fi
}

# Check required environment variables
check_required_vars

# Log configuration
echo "Configuration:"
echo "- GRPC Port: $RANGE_ALLOCATOR_GRPC_PORT"
echo "- Range Default Size: $RANGE_ALLOCATOR_RANGE_DEFAULT_SIZE"
echo "- Range Min Size: $RANGE_ALLOCATOR_RANGE_MIN_SIZE"
echo "- Range Max Size: $RANGE_ALLOCATOR_RANGE_MAX_SIZE"

# Run migrations unless SKIP_MIGRATIONS is set
if [ -z "${SKIP_MIGRATIONS}" ]; then
  echo "Running database migrations..."
  goose -dir /app/migrations postgres "${RANGE_ALLOCATOR_DATABASE_URL}" up
  if [ $? -ne 0 ]; then
    echo "ERROR: Database migration failed"
    exit 1
  fi
  echo "Migrations completed successfully"
else
  echo "Skipping migrations (SKIP_MIGRATIONS is set)"
fi

# Verify grpcurl is available
if ! command -v grpcurl >/dev/null 2>&1; then
  echo "WARNING: grpcurl not found, health checks may fail"
fi

# Start the application
echo "Starting Range Allocator on port ${RANGE_ALLOCATOR_GRPC_PORT}..."
exec /app/range-allocator
