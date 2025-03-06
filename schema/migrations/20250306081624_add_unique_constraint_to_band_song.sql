-- +goose Up
-- +goose StatementBegin
ALTER TABLE songs ADD CONSTRAINT unique_band_song UNIQUE (band, song);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE songs
DROP CONSTRAINT unique_band_song;

-- +goose StatementEnd
