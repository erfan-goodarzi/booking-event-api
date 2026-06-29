-- +goose Up
-- +goose StatementBegin
ALTER TABLE playlist_events
ADD COLUMN IF NOT EXISTS public BOOLEAN;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE playlist_events DROP COLUMN IF EXISTS public;
-- +goose StatementEnd


