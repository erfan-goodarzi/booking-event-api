-- +goose up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS events (
  id BIGSERIAL PRIMARY KEY,
	title VARCHAR(50) NOT NULL,
	location VARCHAR(50) NOT NULL,
	date_time TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	user_id BIGINT NOT NULL REFERENCES users(id),
	description text,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
-- +goose StatementEnd
