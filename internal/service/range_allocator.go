package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/RajNykDhulapkar/gotiny-range-allocator/internal/config"
	"github.com/RajNykDhulapkar/gotiny-range-allocator/internal/repository"
	"github.com/RajNykDhulapkar/gotiny-range-allocator/pkg/db"
	"github.com/google/uuid"
)

type RangeAllocatorService interface {
	AllocateRange(ctx context.Context, params AllocateRangeParams) (*db.Range, error)
	GetRange(ctx context.Context, rangeID uuid.UUID) (*db.Range, error)
	UpdateRangeStatus(ctx context.Context, params UpdateRangeStatusParams) (*db.Range, error)
	GetHealth(ctx context.Context) (bool, string)
}

// RangeAllocator implements the core business logic for range allocation
type RangeAllocator struct {
	repo   repository.RepositoryService
	config *config.RangeConfig
}

// NewRangeAllocator creates a new RangeAllocator service
func NewRangeAllocator(repo repository.RepositoryService, cfg *config.RangeConfig) RangeAllocatorService {
	return &RangeAllocator{
		repo:   repo,
		config: cfg,
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
	size := s.config.DefaultSize
	if params.Size != nil {
		size = *params.Size
		if size < s.config.MinSize || size > s.config.MaxSize {
			return nil, fmt.Errorf(
				"range size must be between %d and %d",
				s.config.MinSize,
				s.config.MaxSize,
			)
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
		Valid:  true,
	}

	// Create range allocation request
	createParams := db.CreateRangeParams{
		ServiceID: params.ServiceID,
		Region:    pgRegion,
		Status:    db.RangeStatusACTIVE,
	}

	return s.repo.AllocateRange(ctx, createParams, size)
}

// GetRange retrieves a range by ID
func (s *RangeAllocator) GetRange(ctx context.Context, rangeID uuid.UUID) (*db.Range, error) {
	return s.repo.GetRange(ctx, rangeID)
}

// UpdateRangeStatusParams contains parameters for updating range status
type UpdateRangeStatusParams struct {
	RangeID   uuid.UUID
	ServiceID string
	Status    db.RangeStatus
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
