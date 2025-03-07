package config

import (
	"strconv"
)

type ServicesConfig struct {
	InfoServiceURL string
	Timeout        int
}

func (c *Config) loadServicesConfig() {

	timeout, err := strconv.Atoi(c.getEnv("EMLIB_INFOSERVICE_TIMEOUT", "500"))
	if err != nil {
		c.Logger.Error("Error: EMLIB_INFOSERVICE_TIMEOUT must be an integer")
	}
	c.Services = ServicesConfig{
		InfoServiceURL: c.getEnv("EMLIB_INFOSERVICE_URL", "http://127.0.0.1:8000"),
		Timeout:        timeout,
	}
}
