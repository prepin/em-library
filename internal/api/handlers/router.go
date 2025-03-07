package handlers

import (
	"em-library/config"
	"em-library/internal/usecase"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	config *config.Config
	Songs  *SongsHandler
	Lyrics *LyricsHandler
}

func NewHandlers(cfg *config.Config, usecases usecase.UseCases) *Handlers {
	return &Handlers{
		config: cfg,
		Songs:  NewSongsHandler(cfg.Logger, usecases),
		Lyrics: NewLyricsHandler(cfg.Logger, usecases),
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

			// Песни
			g.GET("/songs", h.Songs.GetSongsList)
			g.POST("/song", h.Songs.CreateSong)
			g.PATCH("/song/:id", h.Songs.UpdateSong)
			g.DELETE("/song/:id", h.Songs.DeleteSong)

			// Тексты
			g.GET("/song/:id/lyrics", h.Lyrics.GetLyrics)
		}
	}

	registerRoutes(api, apiV1)
}
