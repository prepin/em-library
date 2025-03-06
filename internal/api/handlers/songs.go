package handlers

import (
	"em-library/config"
	"em-library/internal/entities"
	"em-library/internal/errs"
	"em-library/internal/usecase"
	"errors"
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
	Band string `json:"group" binding:"required,min=1"`
	Song string `json:"song" binding:"required,min=1"`
}

type SongResponse struct {
	ID           int    `json:"id"`
	Band         string `json:"group"` // в API инфосервиса используется термин group, сделаем единообразно
	Song         string `json:"song"`
	ReleasedDate string `json:"released_date"`
	Link         string `json:"link"`
}

func (h *SongsHandler) CreateSong(c *gin.Context) {
	var params *CreateSongParams

	if err := c.ShouldBindJSON(&params); err != nil {
		h.logger.Debug("Failed parsing credit request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: formatValidationError(err)})
		return
	}

	songData, err := h.usecases.CreateSong.Execute(
		c.Request.Context(), entities.NewSongData{
			Song: params.Song,
			Band: params.Band,
		},
	)
	if err != nil {
		if errors.Is(err, errs.ErrAlreadyExists) {
			h.logger.Debug("Song already exists", "error", err)
			c.JSON(http.StatusConflict, AlreadyExistsResponse)
			return
		}

		h.logger.Error("Creation of song failed", "error", err)
		c.JSON(http.StatusInternalServerError, ServerErrorResponse)
		return
	}

	h.logger.Info("Song created successfully", "id", songData.ID)

	c.JSON(http.StatusCreated, SongResponse{
		ID:           songData.ID,
		Band:         songData.Band,
		Song:         songData.Song,
		ReleasedDate: songData.ReleaseDate.Format("02.01.2006"),
		Link:         songData.Link,
	})

}
