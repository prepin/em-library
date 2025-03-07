package handlers

import (
	"em-library/config"
	"em-library/internal/usecase"

	docs "em-library/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	// Swagger документация
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Title = "EM Song Library API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Description = "Библиотека песен"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

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
