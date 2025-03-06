package handlers_test

import (
	"em-library/internal/api/handlers"
	"em-library/internal/entities"
	"em-library/internal/errs"
	"em-library/internal/usecase"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupGetLyricsRouter(mockLogger *MockLogger, mockUseCases *MockGetSongLyricsUseCase) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	useCases := usecase.UseCases{
		GetSongLyrics: mockUseCases,
	}

	handler := handlers.NewLyricsHandler(mockLogger, useCases)
	r.GET("/songs/:id/lyrics", handler.GetLyrics)
	return r
}

// Успешное получение текста песни
func TestLyricsHandler_GetLyrics_Success(t *testing.T) {
	mockLogger := new(MockLogger)
	mockUseCase := new(MockGetSongLyricsUseCase)

	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Error", mock.Anything, mock.Anything).Maybe()

	songID := 123
	expectedLyrics := []entities.LyricsVerseData{
		{Index: 1, Content: "This is verse one"},
		{Index: 2, Content: "This is verse two"},
	}

	mockUseCase.On("Execute", mock.Anything, songID, entities.LyricsFilterData{
		Offset: nil,
		Limit:  nil,
	}).Return(expectedLyrics, nil)

	router := setupGetLyricsRouter(mockLogger, mockUseCase)

	req, _ := http.NewRequest(http.MethodGet, "/songs/123/lyrics", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response []entities.LyricsVerseData
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, len(expectedLyrics), len(response))
	assert.Equal(t, expectedLyrics[0].Index, response[0].Index)
	assert.Equal(t, expectedLyrics[0].Content, response[0].Content)
	assert.Equal(t, expectedLyrics[1].Index, response[1].Index)
	assert.Equal(t, expectedLyrics[1].Content, response[1].Content)

	mockUseCase.AssertExpectations(t)
}

// Успешное получение текста песни с параметрами offset и limit
func TestLyricsHandler_GetLyrics_WithPagination(t *testing.T) {
	mockLogger := new(MockLogger)
	mockUseCase := new(MockGetSongLyricsUseCase)

	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Error", mock.Anything, mock.Anything).Maybe()

	songID := 123
	offset := 1
	limit := 10
	expectedLyrics := []entities.LyricsVerseData{
		{Index: 2, Content: "This is verse two"},
	}

	mockUseCase.On("Execute", mock.Anything, songID, entities.LyricsFilterData{
		Offset: &offset,
		Limit:  &limit,
	}).Return(expectedLyrics, nil)

	router := setupGetLyricsRouter(mockLogger, mockUseCase)

	req, _ := http.NewRequest(http.MethodGet, "/songs/123/lyrics?offset=1&limit=10", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response []entities.LyricsVerseData
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, len(expectedLyrics), len(response))
	assert.Equal(t, expectedLyrics[0].Index, response[0].Index)
	assert.Equal(t, expectedLyrics[0].Content, response[0].Content)

	mockUseCase.AssertExpectations(t)
}

// Неверный формат ID песни
func TestLyricsHandler_GetLyrics_InvalidID(t *testing.T) {
	mockLogger := new(MockLogger)
	mockUseCase := new(MockGetSongLyricsUseCase)

	mockLogger.On("Debug", "Missing or invalid ID param for request", mock.Anything).Once()

	router := setupGetLyricsRouter(mockLogger, mockUseCase)

	req, _ := http.NewRequest(http.MethodGet, "/songs/invalid/lyrics", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	mockLogger.AssertExpectations(t)
}

// Неверные параметры пагинации
func TestLyricsHandler_GetLyrics_InvalidPaginationParams(t *testing.T) {
	mockLogger := new(MockLogger)
	mockUseCase := new(MockGetSongLyricsUseCase)

	mockLogger.On("Debug", "Failed parsing request params", mock.Anything).Once()

	router := setupGetLyricsRouter(mockLogger, mockUseCase)

	req, _ := http.NewRequest(http.MethodGet, "/songs/123/lyrics?limit=0", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	mockLogger.AssertExpectations(t)
}

// Песня не найдена
func TestLyricsHandler_GetLyrics_SongNotFound(t *testing.T) {
	mockLogger := new(MockLogger)
	mockUseCase := new(MockGetSongLyricsUseCase)

	mockLogger.On("Debug", "No lyrics found", mock.Anything).Once()

	songID := 999
	mockUseCase.On("Execute", mock.Anything, songID, entities.LyricsFilterData{
		Offset: nil,
		Limit:  nil,
	}).Return(nil, errs.ErrNotFound)

	router := setupGetLyricsRouter(mockLogger, mockUseCase)

	req, _ := http.NewRequest(http.MethodGet, "/songs/999/lyrics", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusNotFound, recorder.Code)

	mockLogger.AssertExpectations(t)
	mockUseCase.AssertExpectations(t)
}

// Внутренняя ошибка сервера при получении текста песни
func TestLyricsHandler_GetLyrics_InternalError(t *testing.T) {
	mockLogger := new(MockLogger)
	mockUseCase := new(MockGetSongLyricsUseCase)

	mockLogger.On("Error", "Getting song lyrics failed", mock.Anything).Once()

	songID := 123
	mockUseCase.On("Execute", mock.Anything, songID, entities.LyricsFilterData{
		Offset: nil,
		Limit:  nil,
	}).Return(nil, errors.New("database error"))

	router := setupGetLyricsRouter(mockLogger, mockUseCase)

	req, _ := http.NewRequest(http.MethodGet, "/songs/123/lyrics", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	mockLogger.AssertExpectations(t)
	mockUseCase.AssertExpectations(t)
}
