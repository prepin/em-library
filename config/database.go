package config

import (
	"fmt"
	"strconv"
)

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func (c *Config) loadDBConfig() {
	port, err := strconv.Atoi(getEnv("EMLIB_DB_PORT", "5432"))
	if err != nil {
		c.Logger.Error("Error: EMLIB must be an integer")
	}

	c.DB = DBConfig{
		Host:     getEnv("EMLIB_DB_HOST", "localhost"),
		Port:     port,
		User:     getEnv("EMLIB_DB_USER", "postgres"),
		Password: getEnv("EMLIB_DB_PASSWORD", ""),
		DBName:   getEnv("EMLIB_NAME", "postgres"),
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
