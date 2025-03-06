package usecase

type UseCases struct {
	CreateSong CreateSongUseCase
}

func NewUseCases(r Repos, s Services) UseCases {
	return UseCases{
		CreateSong: NewCreateSongUseCase(r.TransactionManager, r.SongRepo, r.LyricsRepo, s.SongInfoService),
	}
}
