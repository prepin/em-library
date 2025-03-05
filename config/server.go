package config

import (
	"strconv"
)

type ServerConfig struct {
	Port           string
	ReadTimeout    int
	WriteTimeout   int
	ProductionMode bool
}

func (c *Config) loadServerConfig() {
	readTimeout, err := strconv.Atoi(getEnv("EMLIB_SERVER_READ_TIMEOUT", "5"))
	if err != nil {
		c.Logger.Error("Error: EMLIB_SERVER_READ_TIMEOUT must be an integer")
	}

	writeTimeout, err := strconv.Atoi(getEnv("AV_SERVER_WRITE_TIMEOUT", "5"))
	if err != nil {
		c.Logger.Error("Error: EMLIB_SERVER_WRITE_TIMEOUT must be an integer")
	}

	var productionMode bool
	mode := getEnv("EMLIB_SERVER_MODE", "debug")
	if mode != "debug" && mode != "production" {
		c.Logger.Error("Error: EMLIB_SERVER_MODE must be \"debug\" or \"production\". Setting debug level.")
	}
	if mode == "production" {
		productionMode = true
	}

	c.Server = ServerConfig{
		Port:           getEnv("EMLIB_SERVER_PORT", ":8080"),
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		ProductionMode: productionMode,
	}
}
