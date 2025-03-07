package app

import (
	"em-library/config"
	"em-library/internal/api/handlers"
	"em-library/internal/repository"
	"em-library/internal/services"
	"em-library/internal/usecase"
	"em-library/pkg/database"
)

type Application struct {
	Handlers *handlers.Handlers
}

func New(cfg *config.Config, db *database.Database) *Application {
	repos := usecase.Repos{
		TransactionManager: db.TransactionManager,
		SongRepo:           repository.NewPGSongRepository(db, cfg.Logger),
		LyricsRepo:         repository.NewPGLyricsRepository(db, cfg.Logger),
	}

	services := usecase.Services{
		SongInfoService: services.NewRESTSongInfoService(cfg.Services, cfg.Logger),
	}

	usecases := usecase.NewUseCases(repos, services)

	handlers := handlers.NewHandlers(cfg, usecases)

	return &Application{
		Handlers: handlers,
	}

}
