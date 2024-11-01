// internal/repository/repository.go
package repository

import (
	"context"
	"fmt"

	"github.com/RajNykDhulapkar/gotiny-range-allocator/pkg/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	queries *db.Queries
	pool    *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{
		queries: db.New(pool),
		pool:    pool,
	}
}

func (r *Repository) AllocateRange(ctx context.Context, params db.CreateRangeParams) (*db.Range, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	// Get last range for this service
	lastRange, err := qtx.GetLastRangeForService(ctx, db.GetLastRangeForServiceParams{
		ServiceID: params.ServiceID,
		Region:    params.Region,
	})

	startID := int64(1) // Default start for new service
	if err == nil {
		startID = lastRange.EndID + 1
	} else if err != pgx.ErrNoRows {
		return nil, fmt.Errorf("get last range: %w", err)
	}

	// Create new range
	params.StartID = startID
	params.EndID = startID + 999 // Allocate 1000 IDs by default

	newRange, err := qtx.CreateRange(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("create range: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return newRange, nil
}

func (r *Repository) GetRange(ctx context.Context, rangeID uuid.UUID) (*db.Range, error) {
	rng, err := r.queries.GetRange(ctx, rangeID)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("range not found: %s", rangeID)
	}
	if err != nil {
		return nil, fmt.Errorf("get range: %w", err)
	}
	return rng, nil
}

func (r *Repository) UpdateRangeStatus(ctx context.Context, params db.UpdateRangeStatusParams) (*db.Range, error) {
	rng, err := r.queries.UpdateRangeStatus(ctx, params)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("range not found or service_id mismatch: %s", params.RangeID)
	}
	if err != nil {
		return nil, fmt.Errorf("update range status: %w", err)
	}
	return rng, nil
}

func (r *Repository) CountRanges(ctx context.Context, serviceID string) (int64, error) {
	count, err := r.queries.CountRanges(ctx, db.CountRangesParams{
		ServiceID: serviceID,
	})
	if err != nil {
		return 0, fmt.Errorf("count ranges: %w", err)
	}
	return count, nil
}

func (r *Repository) GetRangesByStatus(ctx context.Context, params db.GetRangesByStatusParams) ([]*db.Range, error) {
	ranges, err := r.queries.GetRangesByStatus(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("get ranges by status: %w", err)
	}
	return ranges, nil
}

func (r *Repository) GetServiceRanges(ctx context.Context, params db.GetServiceRangesParams) ([]*db.Range, error) {
	ranges, err := r.queries.GetServiceRanges(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("get service ranges: %w", err)
	}
	return ranges, nil
}
