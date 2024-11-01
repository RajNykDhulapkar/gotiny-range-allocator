package adapters

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/RajNykDhulapkar/gotiny-range-allocator/internal/service"
	"github.com/RajNykDhulapkar/gotiny-range-allocator/pkg/pb"
	"github.com/RajNykDhulapkar/gotiny-range-allocator/pkg/util"
	"github.com/google/uuid"
)

type grpcAdapter struct {
	pb.UnimplementedRangeAllocatorServer
	rangeAllocator service.RangeAllocatorService
}

func NewGRPCAdapter(allocator service.RangeAllocatorService) pb.RangeAllocatorServer {
	return &grpcAdapter{
		rangeAllocator: allocator,
	}
}

func (h *grpcAdapter) AllocateRange(ctx context.Context, req *pb.AllocateRangeRequest) (*pb.AllocateRangeResponse, error) {
	if req.ServiceId == "" {
		return nil, status.Error(codes.InvalidArgument, "service_id is required")
	}

	params := service.AllocateRangeParams{
		ServiceID: req.ServiceId,
		Size:      req.Size,
	}

	if req.Region != nil && *req.Region != "" {
		params.Region = req.Region
	}

	rng, err := h.rangeAllocator.AllocateRange(ctx, params)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to allocate range: %v", err)
	}

	searializedRange, err := util.ConvertRangeToProto(rng)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to allocate range: %v", err)
	}

	return &pb.AllocateRangeResponse{
		Range: searializedRange,
	}, nil
}

func (h *grpcAdapter) GetRange(ctx context.Context, req *pb.GetRangeRequest) (*pb.Range, error) {
	if req.RangeId == "" {
		return nil, status.Error(codes.InvalidArgument, "range_id is required")
	}

	rangeID, err := uuid.Parse(req.RangeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid range_id format: %v", err)
	}

	rng, err := h.rangeAllocator.GetRange(ctx, rangeID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get range: %v", err)
	}

	searializedRange, err := util.ConvertRangeToProto(rng)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to allocate range: %v", err)
	}

	return searializedRange, nil
}

func (h *grpcAdapter) UpdateRangeStatus(ctx context.Context, req *pb.UpdateRangeStatusRequest) (*pb.Range, error) {
	if req.RangeId == "" {
		return nil, status.Error(codes.InvalidArgument, "range_id is required")
	}
	if req.ServiceId == "" {
		return nil, status.Error(codes.InvalidArgument, "service_id is required")
	}

	rangeID, err := uuid.Parse(req.RangeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid range_id format: %v", err)
	}

	rangeStatus, err := util.ConvertPBToDBRangeStatus(&req.Status)
	if err != nil {
		return nil, fmt.Errorf("invalid RangeStatus value: %v", err)
	}

	params := service.UpdateRangeStatusParams{
		RangeID:   rangeID,
		ServiceID: req.ServiceId,
		Status:    rangeStatus,
	}

	rng, err := h.rangeAllocator.UpdateRangeStatus(ctx, params)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update range status: %v", err)
	}

	searializedRange, err := util.ConvertRangeToProto(rng)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to allocate range: %v", err)
	}

	return searializedRange, nil
}

func (h *grpcAdapter) GetHealth(ctx context.Context, _ *emptypb.Empty) (*pb.HealthResponse, error) {
	isHealthy, details := h.rangeAllocator.GetHealth(ctx)

	status := pb.ServiceStatus_SERVICE_STATUS_NOT_SERVING
	if isHealthy {
		status = pb.ServiceStatus_SERVICE_STATUS_SERVING
	}

	return &pb.HealthResponse{
		Status:  status,
		Details: details,
	}, nil
}
