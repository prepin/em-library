package config

import (
	"fmt"
	"strconv"
)

type DBConfig struct {
	Host          string
	Port          int
	User          string
	Password      string
	DBName        string
	RunMigrations bool
}

func (c *Config) loadDBConfig() {
	port, err := strconv.Atoi(c.getEnv("EMLIB_DB_PORT", "5432"))
	if err != nil {
		c.Logger.Error("Error: EMLIB_DB_PORT must be an integer")
	}

	runMigrationsStr := c.getEnv("EMLIB_RUN_MIGRATIONS", "1")

	c.DB = DBConfig{
		Host:          c.getEnv("EMLIB_DB_HOST", "localhost"),
		Port:          port,
		User:          c.getEnv("EMLIB_DB_USER", "postgres"),
		Password:      c.getEnv("EMLIB_DB_PASSWORD", ""),
		DBName:        c.getEnv("EMLIB_DB_NAME", "postgres"),
		RunMigrations: runMigrationsStr == "1",
	}

}

func (config DBConfig) GetConnectionString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
	)
}
