package config

type ServerConfig struct {
	Port           string
	ProductionMode bool
}

func (c *Config) loadServerConfig() {
	var productionMode bool
	mode := c.getEnv("EMLIB_SERVER_MODE", "debug")
	if mode != "debug" && mode != "production" {
		c.Logger.Error("Error: EMLIB_SERVER_MODE must be \"debug\" or \"production\". Setting debug level.")
	}
	if mode == "production" {
		productionMode = true
	}

	c.Server = ServerConfig{
		Port:           c.getEnv("EMLIB_SERVER_PORT", ":8080"),
		ProductionMode: productionMode,
	}
}
