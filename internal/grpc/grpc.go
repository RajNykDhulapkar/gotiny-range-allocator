// internal/adapters/grpc_handler.go

package adapters

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/RajNykDhulapkar/gotiny-range-allocator/pkg/pb"
)

type grpcAdapter struct {
	pb.UnimplementedRangeAllocatorServer
}

func NewGRPCAdapter() pb.RangeAllocatorServer {
	return &grpcAdapter{}
}

func (h *grpcAdapter) AllocateRange(ctx context.Context, req *pb.AllocateRangeRequest) (*pb.AllocateRangeResponse, error) {
	if req.ServiceId == "" {
		return nil, status.Error(codes.InvalidArgument, "service_id is required")
	}

	// Mock response with fake range
	now := time.Now()
	mockRange := &pb.Range{
		RangeId:     "mock-range-123",
		StartId:     1000,
		EndId:       2000,
		ServiceId:   req.ServiceId,
		Region:      "us-west-1",
		Status:      pb.RangeStatus_RANGE_STATUS_ACTIVE,
		AllocatedAt: timestamppb.New(now),
		UpdatedAt:   timestamppb.New(now),
	}

	return &pb.AllocateRangeResponse{
		Range: mockRange,
	}, nil
}

func (h *grpcAdapter) GetRange(ctx context.Context, req *pb.GetRangeRequest) (*pb.Range, error) {
	if req.RangeId == "" {
		return nil, status.Error(codes.InvalidArgument, "range_id is required")
	}

	// Mock response
	now := time.Now()
	return &pb.Range{
		RangeId:     req.RangeId,
		StartId:     1000,
		EndId:       2000,
		ServiceId:   "mock-service",
		Region:      "us-west-1",
		Status:      pb.RangeStatus_RANGE_STATUS_ACTIVE,
		AllocatedAt: timestamppb.New(now.Add(-24 * time.Hour)), // Allocated yesterday
		UpdatedAt:   timestamppb.New(now),
	}, nil
}

func (h *grpcAdapter) ListRanges(ctx context.Context, req *pb.ListRangesRequest) (*pb.ListRangesResponse, error) {
	if req.ServiceId == "" {
		return nil, status.Error(codes.InvalidArgument, "service_id is required")
	}

	// Mock response with multiple ranges
	now := time.Now()
	mockRanges := []*pb.Range{
		{
			RangeId:     "mock-range-1",
			StartId:     1000,
			EndId:       2000,
			ServiceId:   req.ServiceId,
			Region:      "us-west-1",
			Status:      pb.RangeStatus_RANGE_STATUS_ACTIVE,
			AllocatedAt: timestamppb.New(now.Add(-24 * time.Hour)),
			UpdatedAt:   timestamppb.New(now),
		},
		{
			RangeId:     "mock-range-2",
			StartId:     2001,
			EndId:       3000,
			ServiceId:   req.ServiceId,
			Region:      "us-west-1",
			Status:      pb.RangeStatus_RANGE_STATUS_EXHAUSTED,
			AllocatedAt: timestamppb.New(now.Add(-48 * time.Hour)),
			UpdatedAt:   timestamppb.New(now.Add(-1 * time.Hour)),
		},
	}

	return &pb.ListRangesResponse{
		Ranges:        mockRanges,
		NextPageToken: "", // No more pages
		TotalCount:    2,
	}, nil
}

func (h *grpcAdapter) UpdateRangeStatus(ctx context.Context, req *pb.UpdateRangeStatusRequest) (*pb.Range, error) {
	if req.RangeId == "" {
		return nil, status.Error(codes.InvalidArgument, "range_id is required")
	}
	if req.ServiceId == "" {
		return nil, status.Error(codes.InvalidArgument, "service_id is required")
	}

	// Mock updated range
	now := time.Now()
	return &pb.Range{
		RangeId:     req.RangeId,
		StartId:     1000,
		EndId:       2000,
		ServiceId:   req.ServiceId,
		Region:      "us-west-1",
		Status:      req.Status,
		AllocatedAt: timestamppb.New(now.Add(-24 * time.Hour)),
		UpdatedAt:   timestamppb.New(now),
	}, nil
}

func (h *grpcAdapter) GetHealth(ctx context.Context, _ *emptypb.Empty) (*pb.HealthResponse, error) {
	// Mock healthy response
	return &pb.HealthResponse{
		Status:  pb.ServiceStatus_SERVICE_STATUS_SERVING,
		Details: "Mock service is healthy",
	}, nil
}
