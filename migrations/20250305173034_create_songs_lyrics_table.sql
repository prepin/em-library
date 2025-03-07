-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS lyrics (
  song_id INTEGER PRIMARY KEY,
  content TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW (),
  updated_at TIMESTAMP DEFAULT NOW (),
  FOREIGN KEY (song_id) REFERENCES songs (id) ON DELETE CASCADE
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE lyrics;

-- +goose StatementEnd
