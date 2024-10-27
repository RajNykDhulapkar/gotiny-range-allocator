// cmd/main.go

package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	adapters "github.com/RajNykDhulapkar/gotiny-range-allocator/internal/grpc"
	"github.com/RajNykDhulapkar/gotiny-range-allocator/pkg/pb"
)

func main() {
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

	// Create gRPC server
	server := grpc.NewServer()

	// Register handlers
	handler := adapters.NewGRPCAdapter()
	pb.RegisterRangeAllocatorServer(server, handler)

	// Register reflection service for grpcurl
	reflection.Register(server)

	// Start server in a goroutine
	go func() {
		log.Printf("Starting gRPC server on port %s", port)
		if err := server.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Wait for interrupt signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	// Graceful shutdown
	log.Println("Shutting down gRPC server...")
	server.GracefulStop()
	log.Println("Server stopped")
}
