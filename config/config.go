package config

import (
	"em-library/pkg/logging"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Logger Logger
	Server ServerConfig
	DB     DBConfig
}

func Load() *Config {
	config := &Config{}
	config.load()

	return config
}

func (c *Config) load() {
	err := godotenv.Load()
	if err != nil {
		c.Logger.Debug(".env file not found, only ENV variables are used")
	}

	logLevel := getEnv("EMLIB_LOG_LEVEL", "info")
	switch logLevel {
	case "debug", "info", "warning", "error":
		c.Logger = logging.NewLogger(logLevel)
	default:
		c.Logger = logging.NewLogger("info")
	}

	c.loadDBConfig()
	c.loadServerConfig()
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
