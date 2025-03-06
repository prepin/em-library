-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS songs (
  id SERIAL PRIMARY KEY,
  band VARCHAR(500),
  song VARCHAR(500),
  release_date date,
  link VARCHAR(200),
  created_at TIMESTAMP DEFAULT NOW (),
  updated_at TIMESTAMP DEFAULT NOW ()
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE songs;

-- +goose StatementEnd
