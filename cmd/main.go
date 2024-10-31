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

	adapters "github.com/RajNykDhulapkar/gotiny-range-allocator/internal/grpc"
	"github.com/RajNykDhulapkar/gotiny-range-allocator/internal/repository"
	"github.com/RajNykDhulapkar/gotiny-range-allocator/internal/service"
	"github.com/RajNykDhulapkar/gotiny-range-allocator/pkg/pb"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	ctx := context.Background()

	// Initialize database connection
	dbpool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbpool.Close()

	// Initialize repository
	repo := repository.New(dbpool)

	// Initialize range allocator service
	rangeAllocator := service.NewRangeAllocator(repo)

	// Get port from environment or use default
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50051"
	}

	// Create listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
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
		log.Printf("Starting gRPC server on port %s", port)
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
