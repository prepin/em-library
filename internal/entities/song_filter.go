package entities

import "time"

type SongFilterData struct {
	ID              *int
	Band            *string
	Song            *string
	ReleaseDateFrom *time.Time
	ReleaseDateTo   *time.Time
	Offset          *int
	Limit           *int
}
