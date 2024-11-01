// cmd/main.go

package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/RajNykDhulapkar/gotiny-range-allocator/internal/config"
	adapters "github.com/RajNykDhulapkar/gotiny-range-allocator/internal/grpc"
	"github.com/RajNykDhulapkar/gotiny-range-allocator/internal/repository"
	"github.com/RajNykDhulapkar/gotiny-range-allocator/internal/service"
	"github.com/RajNykDhulapkar/gotiny-range-allocator/pkg/pb"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if err := config.ValidateConfig(cfg); err != nil {
		log.Fatalf("Configuration validation error: %v", err)
	}

	log.Printf("Starting with configuration: DefaultSize=%d, MinSize=%d, MaxSize=%d",
		cfg.Range.DefaultSize,
		cfg.Range.MinSize,
		cfg.Range.MaxSize,
	)

	ctx := context.Background()

	// Initialize database connection
	dbpool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbpool.Close()

	// Initialize repository
	repo := repository.New(dbpool)

	// Initialize range allocator service
	rangeAllocator := service.NewRangeAllocator(repo, &cfg.Range)

	// Create listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Initialize gRPC server
	grpcServer := grpc.NewServer()

	// Initialize and register gRPC handler
	handler := adapters.NewGRPCAdapter(rangeAllocator)
	pb.RegisterRangeAllocatorServer(grpcServer, handler)

	// Register reflection service for grpcurl
	reflection.Register(grpcServer)

	// Start server in a goroutine
	go func() {
		log.Printf("Starting gRPC server on port %s", cfg.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Wait for interrupt signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	// Graceful shutdown
	log.Println("Shutting down gRPC server...")
	grpcServer.GracefulStop()
	log.Println("Server stopped")
}
