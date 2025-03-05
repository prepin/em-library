package handlers

import (
	"em-library/config"
	"em-library/internal/entities"
	"em-library/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SongsHandler struct {
	logger   config.Logger
	usecases usecase.UseCases
}

func NewSongsHandler(l config.Logger, u usecase.UseCases) *SongsHandler {
	return &SongsHandler{
		logger:   l,
		usecases: u,
	}
}

type CreateSongParams struct {
	Group string `json:"group" binding:"required,min=1"`
	Song  string `json:"song" binding:"required,min=1"`
}

type SongResponse struct {
	ID           int    `json:"id"`
	Group        string `json:"group"`
	Song         string `json:"song"`
	ReleasedDate string `json:"released_date"`
	Link         string `json:"link"`
}

func (h *SongsHandler) CreateSong(c *gin.Context) {
	var params *CreateSongParams

	// TODO: добавить валидацию полей в входящей структуре

	if err := c.ShouldBindJSON(&params); err != nil {
		h.logger.Debug("Failed parsing credit request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: formatValidationError(err)})
		return
	}

	songData, err := h.usecases.CreateSong.Execute(
		c.Request.Context(), entities.NewSongData{
			Song:  params.Song,
			Group: params.Group,
		},
	)
	if err != nil {
		h.logger.Error("Creation of song failed", "error", err)
		c.JSON(http.StatusInternalServerError, ServerErrorResponse)
		return
	}

	c.JSON(http.StatusCreated, SongResponse{
		ID:           songData.ID,
		Group:        songData.Group,
		Song:         songData.Song,
		ReleasedDate: songData.ReleaseDate.Format("02.01.2006"),
		Link:         songData.Link,
	})

}
