package entities

// DTO для создания нового текста песни
type NewLyricsData struct {
	SongID  int
	Content string
}

// DTO для передачи текста песен
type LyricsData struct {
	SongID  int
	Content string
}

// DTO для передачи куплета песни
type LyricsVerseData struct {
	Index   int
	Content string
}
