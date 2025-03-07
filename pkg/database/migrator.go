package database

// Goose использует sql.Open интерфейс pgx.
import (
	"database/sql"
	"em-library/config"
	"fmt"
	"path/filepath"
	"runtime"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type Migrator struct {
	config config.DBConfig
}

func NewMigrator(cfg config.DBConfig) *Migrator {
	return &Migrator{
		config: cfg,
	}
}

func (mg *Migrator) RunMigrations() error {
	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)
	migrationsPath := filepath.Join(currentDir, "..", "..", "migrations")

	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	db, err := sql.Open("pgx", mg.config.GetConnectionString())
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	goose.SetLogger(goose.NopLogger())

	if err := goose.Up(db, absPath); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
