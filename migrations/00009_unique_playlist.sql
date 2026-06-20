-- +goose Up
-- +goose StatementBegin
ALTER TABLE playlist_events
ADD CONSTRAINT unique_playlist_event
UNIQUE (playlist_id, event_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE playlist_events DROP CONSTRAINT IF EXISTS unique_playlist_event;
-- +goose StatementEnd


