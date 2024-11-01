-- +goose Up
-- +goose StatementBegin
--
-- create enum
CREATE TYPE range_status AS ENUM ('UNSPECIFIED', 'ACTIVE', 'EXHAUSTED', 'RELEASED');

-- alter column
ALTER TABLE ranges
    ALTER COLUMN status TYPE range_status
    USING status::range_status;

-- rm existing constraint
ALTER TABLE ranges
    DROP CONSTRAINT IF EXISTS ranges_status_check;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- alter column
ALTER TABLE ranges
    ALTER COLUMN status TYPE VARCHAR(20)
    USING status::VARCHAR;

-- drop enum
DROP TYPE IF EXISTS range_status;

-- add constraint again
ALTER TABLE ranges
    ADD CONSTRAINT ranges_status_check CHECK (status IN ('ACTIVE', 'EXHAUSTED', 'RELEASED'));

-- +goose StatementEnd
