package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/RajNykDhulapkar/gotiny-range-allocator/internal/repository"
	"github.com/RajNykDhulapkar/gotiny-range-allocator/pkg/db"
	"github.com/google/uuid"
)

const (
	defaultRangeSize = int64(1000)
	maxRangeSize     = int64(10000)
	minRangeSize     = int64(100)
)

// RangeAllocator implements the core business logic for range allocation
type RangeAllocator struct {
	repo *repository.Repository
}

// NewRangeAllocator creates a new RangeAllocator service
func NewRangeAllocator(repo *repository.Repository) *RangeAllocator {
	return &RangeAllocator{
		repo: repo,
	}
}

// AllocateRangeParams contains parameters for range allocation
type AllocateRangeParams struct {
	ServiceID string
	Size      *int64
	Region    *string
}

// AllocateRange allocates a new range for a service
func (s *RangeAllocator) AllocateRange(ctx context.Context, params AllocateRangeParams) (*db.Range, error) {
	// Validate and set defaults
	size := defaultRangeSize
	if params.Size != nil {
		size = *params.Size
		if size < minRangeSize || size > maxRangeSize {
			return nil, fmt.Errorf("range size must be between %d and %d", minRangeSize, maxRangeSize)
		}
	}

	region := "default"
	if params.Region != nil && *params.Region != "" {
		region = *params.Region
	}

	if params.ServiceID == "" {
		return nil, fmt.Errorf("service_id is required")
	}

	pgRegion := pgtype.Text{
		String: region,
		Valid:  true, // Mark it as a valid non-null value
	}

	// Create range allocation request
	createParams := db.CreateRangeParams{
		ServiceID: params.ServiceID,
		Region:    pgRegion,
		Status:    "ACTIVE",
	}

	return s.repo.AllocateRange(ctx, createParams)
}

// GetRange retrieves a range by ID
func (s *RangeAllocator) GetRange(ctx context.Context, rangeID uuid.UUID) (*db.Range, error) {
	return s.repo.GetRange(ctx, rangeID)
}

// ListRangesParams contains parameters for listing ranges
type ListRangesParams struct {
	ServiceID string
	Status    *string
	Region    *string
	PageSize  int32
	PageToken *string
}

// ListRangesResult contains the result of listing ranges
type ListRangesResult struct {
	Ranges        []*db.Range
	NextPageToken string
	TotalCount    int64
}

// ListRanges lists ranges for a service with pagination and filtering
func (s *RangeAllocator) ListRanges(ctx context.Context, params ListRangesParams) (*ListRangesResult, error) {
	if params.ServiceID == "" {
		return nil, fmt.Errorf("service_id is required")
	}

	if params.PageSize <= 0 {
		params.PageSize = 50 // default page size
	}

	listParams := db.ListRangesParams{
		ServiceID:    params.ServiceID,
		StatusFilter: "",
		RegionFilter: "",
		CursorID:     uuid.Nil,
		Limit:        params.PageSize + 1, // fetch one extra to determine if there are more pages
	}

	if params.Status != nil {
		listParams.StatusFilter = *params.Status
	}

	if params.Region != nil {
		listParams.RegionFilter = *params.Region
	}

	if params.PageToken != nil && *params.PageToken != "" {
		cursorID, err := uuid.Parse(*params.PageToken)
		if err != nil {
			return nil, fmt.Errorf("invalid page token: %w", err)
		}
		listParams.CursorID = cursorID
	}

	ranges, err := s.repo.ListRanges(ctx, listParams)
	if err != nil {
		return nil, fmt.Errorf("list ranges: %w", err)
	}

	var nextPageToken string
	if len(ranges) > int(params.PageSize) {
		nextPageToken = ranges[len(ranges)-1].RangeID.String()
		ranges = ranges[:len(ranges)-1]
	}

	totalCount, err := s.repo.CountRanges(ctx, params.ServiceID)
	if err != nil {
		return nil, fmt.Errorf("count ranges: %w", err)
	}

	return &ListRangesResult{
		Ranges:        ranges,
		NextPageToken: nextPageToken,
		TotalCount:    totalCount,
	}, nil
}

// UpdateRangeStatusParams contains parameters for updating range status
type UpdateRangeStatusParams struct {
	RangeID   uuid.UUID
	ServiceID string
	Status    string
}

// UpdateRangeStatus updates the status of a range
func (s *RangeAllocator) UpdateRangeStatus(ctx context.Context, params UpdateRangeStatusParams) (*db.Range, error) {
	if params.RangeID == uuid.Nil {
		return nil, fmt.Errorf("range_id is required")
	}
	if params.ServiceID == "" {
		return nil, fmt.Errorf("service_id is required")
	}
	if params.Status == "" {
		return nil, fmt.Errorf("status is required")
	}

	// Validate status
	switch params.Status {
	case "ACTIVE", "EXHAUSTED", "RELEASED":
		// valid status
	default:
		return nil, fmt.Errorf("invalid status: %s", params.Status)
	}

	updateParams := db.UpdateRangeStatusParams{
		RangeID:   params.RangeID,
		ServiceID: params.ServiceID,
		Status:    params.Status,
	}

	return s.repo.UpdateRangeStatus(ctx, updateParams)
}

// GetHealth checks the health of the service and its dependencies
func (s *RangeAllocator) GetHealth(ctx context.Context) (bool, string) {
	// Simple health check: try to count ranges
	_, err := s.repo.CountRanges(ctx, "health-check")
	if err != nil {
		return false, fmt.Sprintf("database health check failed: %v", err)
	}
	return true, "service is healthy"
}
