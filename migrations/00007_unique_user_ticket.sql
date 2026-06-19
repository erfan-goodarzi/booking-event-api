-- +goose Up
-- +goose StatementBegin
ALTER TABLE bookings
ADD CONSTRAINT unique_user_ticket
UNIQUE (user_id, ticket_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE bookings DROP CONSTRAINT IF EXISTS unique_user_ticket;
-- +goose StatementEnd


