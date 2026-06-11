-- +goose Up
-- +goose StatementBegin
ALTER TABLE events
ADD COLUMN IF NOT EXISTS version INT DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE events DROP COLUMN IF EXISTS version;
-- +goose StatementEnd


