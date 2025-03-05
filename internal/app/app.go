package app

import (
	"em-library/config"
	"em-library/internal/app/api/handlers"
)

type Application struct {
	Handlers *handlers.Handlers
}

func New(cfg *config.Config /* db *database.Database, redis *redis.Redis*/) *Application {
	handlers := handlers.NewHandlers(cfg /*, usecases */)

	return &Application{
		Handlers: handlers,
	}

}
