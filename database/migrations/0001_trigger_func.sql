-- +goose Up

-- Generic trigger function for setting updated_at time on row update for any table with that column
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_updated_at() RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = NOW(); RETURN NEW;
END;
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose Down
DROP FUNCTION set_updated_at();
