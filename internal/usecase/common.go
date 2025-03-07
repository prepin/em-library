package usecase

type UseCases struct {
	CreateSong    CreateSongUseCase
	GetSongList   GetSongListUseCase
	GetSongLyrics GetSongLyricsUseCase
	DeleteSong    DeleteSongUseCase
	UpdateSong    UpdateSongUseCase
}

func NewUseCases(r Repos, s Services) UseCases {
	return UseCases{
		CreateSong:    NewCreateSongUseCase(r.TransactionManager, r.SongRepo, r.LyricsRepo, s.SongInfoService),
		GetSongList:   NewGetSongListUseCase(r.SongRepo),
		GetSongLyrics: NewGetSongLyricsUsecase(r.LyricsRepo),
		DeleteSong:    NewDeleteSongUseCase(r.TransactionManager, r.SongRepo, r.LyricsRepo),
		UpdateSong:    NewUpdateSongUseCase(r.TransactionManager, r.SongRepo, r.LyricsRepo),
	}
}
