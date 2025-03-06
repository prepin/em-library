package handlers

import (
	"em-library/config"
	"em-library/internal/entities"
	"em-library/internal/errs"
	"em-library/internal/usecase"
	"errors"
	"net/http"
	"strconv"
	"time"

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

func (h *SongsHandler) CreateSong(c *gin.Context) {
	var params *CreateSongParams

	if err := c.ShouldBindJSON(&params); err != nil {
		h.logger.Debug("Failed parsing request params", "error", err)
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

	h.logger.Debug("Song created successfully", "id", songData.ID)

	c.JSON(http.StatusCreated, songData)

}

type GetSongsParams struct {
	ID              *int       `form:"id" binding:"omitempty,gt=0"`
	Band            *string    `form:"band" binding:"omitempty,min=1"`
	Song            *string    `form:"song" binding:"omitempty,min=1"`
	ReleaseDateFrom *time.Time `form:"release_date_from" binding:"omitempty" time_format:"2006-01-02"`
	ReleaseDateTo   *time.Time `form:"release_date_to" binding:"omitempty" time_format:"2006-01-02"`
	Offset          *int       `form:"offset" binding:"omitempty"`
	Limit           *int       `form:"limit" binding:"omitempty,min=1"`
}

func (h *SongsHandler) GetSongsList(c *gin.Context) {
	var params GetSongsParams

	if err := c.ShouldBindQuery(&params); err != nil {
		h.logger.Debug("Failed parsing request params", "error", err)
		c.JSON(http.StatusBadRequest, InvalidRequestResponse)
		return
	}

	songs, err := h.usecases.GetSongList.Execute(c.Request.Context(), entities.SongFilterData{
		ID:              params.ID,
		Band:            params.Band,
		Song:            params.Song,
		ReleaseDateFrom: params.ReleaseDateFrom,
		ReleaseDateTo:   params.ReleaseDateTo,
		Offset:          params.Offset,
		Limit:           params.Limit,
	})

	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			h.logger.Debug("No songs found", "error", err)
			c.JSON(http.StatusNotFound, NotFoundResponse)
			return
		}

		h.logger.Error("Getting song list failed", "error", err)
		c.JSON(http.StatusInternalServerError, ServerErrorResponse)
		return
	}

	h.logger.Debug("Songs list retrieved successfully")

	c.JSON(http.StatusOK, songs)
}

func (h *SongsHandler) DeleteSong(c *gin.Context) {
	songIDParam := c.Param("id")
	songID, err := strconv.Atoi(songIDParam)

	if err != nil {
		h.logger.Debug("Missing or invalid ID param for request", "ID param", songIDParam)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "song ID is required"})
		return
	}

	err = h.usecases.DeleteSong.Execute(c.Request.Context(), songID)
	if err != nil {
		h.logger.Debug("Failed to delete song", "ID", songID)
		c.JSON(http.StatusInternalServerError, ServerErrorResponse)
		return
	}

	c.Status(http.StatusNoContent)
}
