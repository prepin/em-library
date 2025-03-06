package usecase

type UseCases struct {
	CreateSong  CreateSongUseCase
	GetSongList GetSongListUseCase
}

func NewUseCases(r Repos, s Services) UseCases {
	return UseCases{
		CreateSong:  NewCreateSongUseCase(r.TransactionManager, r.SongRepo, r.LyricsRepo, s.SongInfoService),
		GetSongList: NewGetSongListUseCase(r.SongRepo),
	}
}
