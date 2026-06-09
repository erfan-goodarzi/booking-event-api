-- +goose Up
-- +goose StatementBegin
CREATE TYPE ticket_type AS ENUM ('vip', 'general');
CREATE TYPE booking_status AS ENUM ('confirmed', 'cancelled', 'pending');

CREATE TABLE IF NOT EXISTS tickets (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
  type ticket_type NOT NULL default 'general',
  price NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
  quantity INT NOT NULL Default 0, 
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS bookings (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  ticket_id UUID NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  status booking_status NOT NULL DEFAULT 'confirmed',
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS bookings;
DROP TABLE IF EXISTS tickets;
DROP TYPE IF EXISTS ticket_type;
DROP TYPE IF EXISTS booking_status;
-- +goose StatementEnd
