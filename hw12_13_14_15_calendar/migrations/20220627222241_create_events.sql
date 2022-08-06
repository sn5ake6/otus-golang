-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS events
(
  id UUID NOT NULL PRIMARY KEY,
  title TEXT NOT NULL,
  begin_at TIMESTAMP NOT NULL,
  end_at TIMESTAMP NOT NULL,
  description TEXT NOT NULL,
  user_id UUID NOT NULL,
  notify_at TIMESTAMP NOT NULL,
  notified_at TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS events;
-- +goose StatementEnd
