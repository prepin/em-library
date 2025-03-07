package handlers

import (
	"em-library/config"
	"em-library/internal/entities"
	"em-library/internal/errs"
	"em-library/internal/usecase"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LyricsHandler struct {
	logger   config.Logger
	usecases usecase.UseCases
}

func NewLyricsHandler(l config.Logger, u usecase.UseCases) *LyricsHandler {
	return &LyricsHandler{
		logger:   l,
		usecases: u,
	}
}

type GetLyricsParams struct {
	Offset *int `form:"offset" binding:"omitempty,min=0"`
	Limit  *int `form:"limit" binding:"omitempty,min=1"`
}

func (h *LyricsHandler) GetLyrics(c *gin.Context) {
	songIDParam := c.Param("id")
	songID, err := strconv.Atoi(songIDParam)

	if err != nil {
		h.logger.Debug("Missing or invalid ID param for request", "ID param", songIDParam)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "song ID is required"})
		return
	}

	var params GetLyricsParams
	if err := c.ShouldBindQuery(&params); err != nil {
		h.logger.Debug("Failed parsing request params", "error", err)
		c.JSON(http.StatusBadRequest, InvalidRequestResponse)
		return
	}

	lyrics, err := h.usecases.GetSongLyrics.Execute(c.Request.Context(), songID, entities.LyricsFilterData{
		Offset: params.Offset,
		Limit:  params.Limit,
	})

	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			h.logger.Debug("No lyrics found", "error", err)
			c.JSON(http.StatusNotFound, NotFoundResponse)
			return
		}

		h.logger.Error("Getting song lyrics failed", "error", err)
		c.JSON(http.StatusInternalServerError, ServerErrorResponse)
		return
	}

	c.JSON(http.StatusOK, lyrics)
}
