-- name: CreateRange :one
-- Creates a new range allocation for a service
INSERT INTO ranges (
    start_id,
    end_id,
    service_id,
    region,
    status
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetRange :one
-- Gets a range by its ID
SELECT * FROM ranges
WHERE range_id = $1;

-- name: GetLastRangeForService :one
-- Gets the last allocated range for a service in a specific region
SELECT * FROM ranges
WHERE service_id = $1 
AND region = $2
ORDER BY end_id DESC
LIMIT 1;

-- name: ListRanges :many
-- Lists ranges for a service with optional filters
SELECT * FROM ranges
WHERE service_id = $1
    AND (sqlc.narg('status_filter')::range_status IS NULL OR status = sqlc.narg('status_filter')::range_status)
    AND (sqlc.narg('region_filter')::text IS NULL OR region = sqlc.narg('region_filter'))
    AND (sqlc.narg('cursor_id')::uuid IS NULL OR range_id > sqlc.narg('cursor_id'))
ORDER BY range_id
LIMIT $2;

-- name: CountRanges :one
-- Counts total ranges for a service
SELECT COUNT(*) 
FROM ranges 
WHERE service_id = $1
    AND (sqlc.narg('status_filter')::range_status IS NULL OR status = sqlc.narg('status_filter')::range_status)
    AND (sqlc.narg('region_filter')::text IS NULL OR region = sqlc.narg('region_filter'));

-- name: UpdateRangeStatus :one
-- Updates the status of a range
UPDATE ranges
SET status = $3
WHERE range_id = $1 
AND service_id = $2
RETURNING *;

-- name: GetRangesByStatus :many
-- Gets ranges by status
SELECT * FROM ranges
WHERE status = $1
ORDER BY range_id
LIMIT $2 OFFSET $3;

-- name: GetServiceRanges :many
-- Gets all ranges for a service in a specific region
SELECT * FROM ranges
WHERE service_id = $1 
AND region = $2
ORDER BY start_id;

-- name: DeleteRange :exec
-- Deletes a range (for testing purposes)
DELETE FROM ranges
WHERE range_id = $1
AND service_id = $2;
