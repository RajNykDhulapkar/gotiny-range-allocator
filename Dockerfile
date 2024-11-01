FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" \
    -o bin/range-allocator ./cmd/main.go

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

FROM alpine:3.19

RUN apk add --no-cache \
    ca-certificates \
    curl \
    tar \
    tzdata && \
    curl -sSL "https://github.com/fullstorydev/grpcurl/releases/download/v1.8.7/grpcurl_1.8.7_linux_x86_64.tar.gz" | tar -xz -C /usr/local/bin

RUN adduser -D -g '' appuser

WORKDIR /app

COPY --from=builder /app/bin/range-allocator .
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /app/migrations ./migrations
COPY scripts/startup.sh .

RUN chown -R appuser:appuser /app && \
    chmod +x /app/range-allocator && \
    chmod +x /app/startup.sh

USER appuser

EXPOSE ${RANGE_ALLOCATOR_GRPC_PORT}

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD grpcurl -plaintext localhost:${RANGE_ALLOCATOR_GRPC_PORT} rangeallocator.v1.RangeAllocator/GetHealth || exit 1

ENTRYPOINT ["/app/startup.sh"]
