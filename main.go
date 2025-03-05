package main

import (
	"em-library/config"
	"em-library/internal/app"
	"em-library/pkg/server"
)

func main() {
	cfg := config.Load()
	// db := database.NewDatabase(cfg.DB)
	// defer db.Close()

	app := app.New(cfg /*, db, redis */)

	cfg.Logger.Info("launched song library service", "config", cfg.Server)

	srv := server.New(cfg, app.Handlers)
	if err := srv.Run(); err != nil {
		cfg.Logger.Error("failed to launch server", "Error", err)
	}
}
