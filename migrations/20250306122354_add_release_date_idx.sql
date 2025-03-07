-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_songs_release_date ON songs (release_date);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_songs_release_date;

-- +goose StatementEnd
