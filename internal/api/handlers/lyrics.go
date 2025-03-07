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

// GetLyrics godoc
// @Summary Получить текст песни
// @Description Получить куплеты песни по ID песни с возможностью пагинации
// @Tags lyrics
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Param offset query int false "С какого куплета начать"
// @Param limit query int false "Сколько куплетов вывести для пагинации"
// @Success 200 {array} entities.LyricsVerseData "Текст песни успешно получен"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 404 {object} ErrorResponse "Текст песни не найден"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /song/{id}/lyrics [get]
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
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
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

	h.logger.Info("Song lyrics retrieved successfully", "song", songID)
	c.JSON(http.StatusOK, lyrics)
}
