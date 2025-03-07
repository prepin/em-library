package config

import (
	"em-library/pkg/logging"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Logger   Logger
	Server   ServerConfig
	DB       DBConfig
	Services ServicesConfig
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

	logLevel := os.Getenv("EMLIB_LOG_LEVEL")
	switch logLevel {
	case "debug", "info", "warning", "error":
		c.Logger = logging.NewLogger(logLevel)
	default:
		c.Logger = logging.NewLogger("info")
	}

	c.loadDBConfig()
	c.loadServerConfig()
	c.loadServicesConfig()
}

func (c *Config) getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		c.Logger.Debug("Environment variable not set. Fallback to default value", "key", key)
		return defaultValue
	}
	return value
}
