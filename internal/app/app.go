package app

import (
	"em-library/config"
	"em-library/internal/api/handlers"
	"em-library/internal/repository"
	"em-library/internal/services"
	"em-library/internal/usecase"
)

type Application struct {
	Handlers *handlers.Handlers
}

func New(cfg *config.Config /* db *database.Database, redis *redis.Redis*/) *Application {
	repos := usecase.Repos{
		SongRepo: repository.NewPGSongRepository(),
	}

	services := usecase.Services{
		SongInfoService: services.NewRESTSongInfoService(),
	}

	usecases := usecase.NewUseCases(repos, services)

	handlers := handlers.NewHandlers(cfg, usecases)

	return &Application{
		Handlers: handlers,
	}

}
