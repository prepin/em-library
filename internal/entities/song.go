package entities

import "time"

// DTO для создания новой песни
type NewSongData struct {
	Group string
	Song  string
}

// DTO для полной информации о песне (без текста)
type SongData struct {
	ID          int
	Group       string
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
