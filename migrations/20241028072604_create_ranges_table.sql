-- +goose Up
-- +goose StatementBegin

-- 1. Create the `ranges` table
CREATE TABLE IF NOT EXISTS ranges (
    range_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    start_id BIGINT NOT NULL,
    end_id BIGINT NOT NULL,
    service_id VARCHAR(255) NOT NULL,
    region VARCHAR(50),
    status VARCHAR(20) NOT NULL,
    allocated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT ranges_start_end_check CHECK (start_id < end_id),
    CONSTRAINT ranges_status_check CHECK (status IN ('ACTIVE', 'EXHAUSTED', 'RELEASED'))
);

-- 2. Create indexes for efficient querying
CREATE INDEX IF NOT EXISTS idx_ranges_service_id ON ranges(service_id);
CREATE INDEX IF NOT EXISTS idx_ranges_status ON ranges(status);
CREATE INDEX IF NOT EXISTS idx_ranges_region ON ranges(region);

-- 3. Create function to update `updated_at` timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 4. Create trigger to automatically update `updated_at`
CREATE TRIGGER update_ranges_updated_at
    BEFORE UPDATE ON ranges
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd

-----------------------------------------------------------------
-----------------------------------------------------------------


-- +goose Down
-- +goose StatementBegin

-- 1. Drop trigger
DROP TRIGGER IF EXISTS update_ranges_updated_at ON ranges;

-- 2. Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- 3. Drop indexes
DROP INDEX IF EXISTS idx_ranges_region;
DROP INDEX IF EXISTS idx_ranges_status;
DROP INDEX IF EXISTS idx_ranges_service_id;

-- 4. Drop table
DROP TABLE IF EXISTS ranges;

-- +goose StatementEnd
