package adapters

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/RajNykDhulapkar/gotiny-range-allocator/internal/service"
	"github.com/RajNykDhulapkar/gotiny-range-allocator/pkg/db"
	"github.com/RajNykDhulapkar/gotiny-range-allocator/pkg/pb"
	"github.com/google/uuid"
)

type grpcAdapter struct {
	pb.UnimplementedRangeAllocatorServer
	rangeAllocator *service.RangeAllocator
}

func NewGRPCAdapter(allocator *service.RangeAllocator) pb.RangeAllocatorServer {
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

	if req.Region != "" {
		params.Region = &req.Region
	}

	rng, err := h.rangeAllocator.AllocateRange(ctx, params)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to allocate range: %v", err)
	}

	return &pb.AllocateRangeResponse{
		Range: convertRangeToProto(rng),
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

	return convertRangeToProto(rng), nil
}

func (h *grpcAdapter) ListRanges(ctx context.Context, req *pb.ListRangesRequest) (*pb.ListRangesResponse, error) {
	if req.ServiceId == "" {
		return nil, status.Error(codes.InvalidArgument, "service_id is required")
	}

	params := service.ListRangesParams{
		ServiceID: req.ServiceId,
		PageSize:  req.PageSize,
	}

	if req.Status != pb.RangeStatus_RANGE_STATUS_UNSPECIFIED {
		status := req.Status.String()
		params.Status = &status
	}

	if req.Region != "" {
		params.Region = &req.Region
	}

	if req.PageToken != "" {
		params.PageToken = &req.PageToken
	}

	result, err := h.rangeAllocator.ListRanges(ctx, params)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list ranges: %v", err)
	}

	ranges := make([]*pb.Range, len(result.Ranges))
	for i, r := range result.Ranges {
		ranges[i] = convertRangeToProto(&r)
	}

	return &pb.ListRangesResponse{
		Ranges:        ranges,
		NextPageToken: result.NextPageToken,
		TotalCount:    int32(result.TotalCount),
	}, nil
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

	params := service.UpdateRangeStatusParams{
		RangeID:   rangeID,
		ServiceID: req.ServiceId,
		Status:    req.Status.String(),
	}

	rng, err := h.rangeAllocator.UpdateRangeStatus(ctx, params)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update range status: %v", err)
	}

	return convertRangeToProto(rng), nil
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

func convertRangeToProto(r *db.Range) *pb.Range {
	var status pb.RangeStatus
	switch r.Status {
	case "ACTIVE":
		status = pb.RangeStatus_RANGE_STATUS_ACTIVE
	case "EXHAUSTED":
		status = pb.RangeStatus_RANGE_STATUS_EXHAUSTED
	case "RELEASED":
		status = pb.RangeStatus_RANGE_STATUS_RELEASED
	default:
		status = pb.RangeStatus_RANGE_STATUS_UNSPECIFIED
	}

	return &pb.Range{
		RangeId:     r.RangeID.String(),
		StartId:     r.StartID,
		EndId:       r.EndID,
		ServiceId:   r.ServiceID,
		Region:      r.Region,
		Status:      status,
		AllocatedAt: timestamppb.New(r.AllocatedAt),
		UpdatedAt:   timestamppb.New(r.UpdatedAt),
	}
}
