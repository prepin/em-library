package handlers

import (
	"em-library/config"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	config *config.Config
}

func NewHandlers(cfg *config.Config /*, usecases usecase.Usecases */) *Handlers {
	return &Handlers{
		config: cfg,
	}
}

func (h *Handlers) RegisterRoutes(r *gin.Engine /*, jwtService *auth.JWTService */) {

	api := r.Group("/api")
	apiV1 := r.Group("/api/v1")

	registerRoutes := func(groups ...*gin.RouterGroup) {
		for _, g := range groups {
			// healthcheck эндпойнты
			g.GET("/ping", GetPingHandler)
			g.GET("/teapot", GetTeapotHandler)
		}
	}

	registerRoutes(api, apiV1)
}
