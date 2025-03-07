package entities

import (
	"encoding/json"
	"time"
)

// DTO для создания новой песни
type NewSongData struct {
	Band        string
	Song        string
	ReleaseDate time.Time
	Link        string
}

// DTO для полной информации о песне (без текста)
type SongData struct {
	ID          int       `json:"id"`
	Band        string    `json:"group"`
	Song        string    `json:"song"`
	ReleaseDate time.Time `json:"release_date"`
	Link        string    `json:"link"`
}

func (s SongData) MarshalJSON() ([]byte, error) {
	type Alias SongData
	return json.Marshal(&struct {
		ReleaseDate string `json:"release_date"`
		*Alias
	}{
		ReleaseDate: s.ReleaseDate.Format("2006-01-02"),
		Alias:       (*Alias)(&s),
	})
}

// DTO для обогащения данных песни
type SongDetail struct {
	ReleaseDate time.Time
	Lyrics      string
	Link        string
}

// DTO для обновления данных песни
type UpdateSongData struct {
	Band        *string
	Song        *string
	ReleaseDate *time.Time
	Link        *string
	Lyrics      *string
}
