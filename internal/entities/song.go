package entities

import "time"

// DTO для создания новой песни
type NewSongData struct {
	Band        string
	Song        string
	ReleaseDate time.Time
	Link        string
}

// DTO для полной информации о песне (без текста)
type SongData struct {
	ID          int
	Band        string
	Song        string
	ReleaseDate time.Time
	Link        string
}

// DTO для обогащения данных песни
type SongDetail struct {
	ReleaseDate time.Time
	Text        string
	Link        string
}
