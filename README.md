# GoTiny Range Allocator

A distributed range allocation service built with Go, designed to provide unique ID ranges for multiple URL shortener instances. This service ensures consistent and non-overlapping ID ranges across distributed systems.

## Features

- Distributed range allocation with PostgreSQL backend
- gRPC API for service communication
- Efficient range management with customizable sizes
- Health monitoring and status tracking
- Database migrations using Goose
- Docker containerization
- Automatic schema generation with SQLC

## Architecture

```mermaid
graph TD
    A[URL Shortener Services] -->|gRPC| B[Range Allocator Service]
    B -->|Manages Ranges| C[PostgreSQL]
```

## Tech Stack

- Go 1.21+
- PostgreSQL (UUID, Timestamps)
- Protocol Buffers & gRPC
- SQLC for type-safe SQL
- Goose for migrations
- Docker & Docker Compose
- pgx for PostgreSQL driver

## Quick Start

1. Clone the repository:

```bash
git clone https://github.com/RajNykDhulapkar/gotiny-range-allocator.git
cd gotiny-range-allocator
```

2. Set up development environment:

```bash
./scripts/setup.sh
```

3. Run with Docker:

```bash
docker-compose up -d
```

## Environment Variables

```env
RANGE_ALLOCATOR_GRPC_PORT=50051
RANGE_ALLOCATOR_DATABASE_URL=postgres://username:password@localhost:5432/alloc?sslmode=disable
RANGE_ALLOCATOR_RANGE_DEFAULT_SIZE=999
RANGE_ALLOCATOR_RANGE_MIN_SIZE=100
RANGE_ALLOCATOR_RANGE_MAX_SIZE=10000
```

## API

### gRPC Service Definition

```protobuf
service RangeAllocator {
  rpc AllocateRange(AllocateRangeRequest) returns (AllocateRangeResponse);
  rpc GetRange(GetRangeRequest) returns (Range);
  rpc UpdateRangeStatus(UpdateRangeStatusRequest) returns (Range);
  rpc GetHealth(Empty) returns (HealthResponse);
}
```

## Development

```bash
# Generate Protocol Buffers
make proto-gen

# Create new migration
make migrate-create name=migration_name

# Run migrations
make migrate-up

# Build
make build

# Run locally
make run
```

## Testing

```bash
# Run integration tests
./tests/test_endpoints.sh
```

## Docker Support

Build the image:

```bash
docker build -t range-allocator .
```

Run container:

```bash
docker run -d \
  --name range-allocator \
  -p 50051:50051 \
  -e RANGE_ALLOCATOR_DATABASE_URL="postgres://username:password@postgres:5432/alloc?sslmode=disable" \
  -e RANGE_ALLOCATOR_GRPC_PORT=50051 \
  -e RANGE_ALLOCATOR_RANGE_DEFAULT_SIZE=1000 \
  range-allocator
```

## Project Structure

```
.
├── cmd/                # Application entrypoint
├── internal/
│   ├── config/        # Configuration management
│   ├── grpc/          # gRPC service implementations
│   ├── repository/    # Data access layer
│   └── service/       # Business logic
├── migrations/        # Database migrations
├── pkg/
│   ├── db/           # Generated database code
│   ├── pb/           # Generated Protocol Buffers
│   └── util/         # Shared utilities
└── scripts/          # Development and deployment scripts
```

## Database Schema

The service uses a PostgreSQL database with the following schema:

```sql
CREATE TABLE ranges (
    range_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    start_id BIGINT NOT NULL,
    end_id BIGINT NOT NULL,
    service_id VARCHAR(255) NOT NULL,
    region VARCHAR(50),
    status range_status NOT NULL,
    allocated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
```

## License

MIT License

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request
