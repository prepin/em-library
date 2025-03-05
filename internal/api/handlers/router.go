package handlers

import (
	"em-library/config"
	"em-library/internal/usecase"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	config *config.Config
	Songs  *SongsHandler
}

func NewHandlers(cfg *config.Config, usecases usecase.UseCases) *Handlers {
	return &Handlers{
		config: cfg,
		Songs:  NewSongsHandler(cfg.Logger, usecases),
	}
}

func (h *Handlers) RegisterRoutes(r *gin.Engine) {

	// две группы нужны для версионирования API при возможных изменениях без обратной совместимости
	api := r.Group("/api")
	apiV1 := r.Group("/api/v1")

	registerRoutes := func(groups ...*gin.RouterGroup) {
		for _, g := range groups {
			// healthcheck эндпойнты
			g.GET("/ping", GetPingHandler)
			g.GET("/teapot", GetTeapotHandler)

			// Songs
			g.POST("/song", h.Songs.CreateSong)
		}
	}

	registerRoutes(api, apiV1)
}
