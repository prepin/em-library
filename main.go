package main

import (
	"em-library/config"
	"em-library/internal/app"
	"em-library/pkg/database"
	"em-library/pkg/server"
	"os"
)

func main() {
	cfg := config.Load()

	cfg.Logger.Debug("Running migrations")
	mg := database.NewMigrator(cfg.DB)
	if err := mg.RunMigrations(); err != nil {
		cfg.Logger.Error("Failed to run migrations", "errror", err)
		os.Exit(1)
	}
	cfg.Logger.Debug("Migrations run successfully")

	cfg.Logger.Debug("Connecting to database")
	db := database.NewDatabase(cfg.DB, cfg.Logger)
	if db == nil {
		cfg.Logger.Error("Failed connect to database", "error")
		os.Exit(1)
	}
	defer db.Close()
	cfg.Logger.Debug("Connected to database successfully")

	app := app.New(cfg, db)

	cfg.Logger.Info("launched song library service", "config", cfg.Server)

	srv := server.New(cfg, app.Handlers)
	if err := srv.Run(); err != nil {
		cfg.Logger.Error("failed to launch server", "Error", err)
	}
}
