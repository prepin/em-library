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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupGetSongsRouter(mockLogger *MockLogger, mockUseCases *MockGetSongListUseCase) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	useCases := usecase.UseCases{
		GetSongList: mockUseCases,
	}

	handler := handlers.NewSongsHandler(mockLogger, useCases)
	r.GET("/songs", handler.GetSongsList)
	return r
}

type SongResponse struct {
	ID          int    `json:"id"`
	Band        string `json:"band"`
	Song        string `json:"song"`
	ReleaseDate string `json:"release_date"`
	Link        string `json:"link"`
}

// Успешное получение списка песен без фильтров
func TestSongsHandler_GetSongsList_Success(t *testing.T) {
	mockLogger := new(MockLogger)
	mockUseCase := new(MockGetSongListUseCase)

	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()

	releaseDate := time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC)
	expectedSongs := []entities.SongData{
		{
			ID:          123,
			Band:        "Test Group",
			Song:        "Test Song",
			ReleaseDate: releaseDate,
			Link:        "https://example.com/song1",
		},
		{
			ID:          124,
			Band:        "Another Group",
			Song:        "Another Song",
			ReleaseDate: releaseDate,
			Link:        "https://example.com/song2",
		},
	}

	mockUseCase.On("Execute", mock.Anything, mock.Anything).Return(expectedSongs, nil)

	router := setupGetSongsRouter(mockLogger, mockUseCase)

	req, _ := http.NewRequest(http.MethodGet, "/songs", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response []SongResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, expectedSongs[0].ID, response[0].ID)
	assert.Equal(t, expectedSongs[0].Band, response[0].Band)
	assert.Equal(t, expectedSongs[0].Song, response[0].Song)
	assert.Equal(t, "2022-02-01", response[0].ReleaseDate)
	assert.Equal(t, expectedSongs[0].Link, response[0].Link)

	mockUseCase.AssertExpectations(t)
}

// Успешное получение списка песен с фильтрами
func TestSongsHandler_GetSongsList_WithFilters(t *testing.T) {
	mockLogger := new(MockLogger)
	mockUseCase := new(MockGetSongListUseCase)

	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()

	releaseDate := time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC)
	expectedSongs := []entities.SongData{
		{
			ID:          123,
			Band:        "Test Group",
			Song:        "Test Song",
			ReleaseDate: releaseDate,
			Link:        "https://example.com/song1",
		},
	}

	mockUseCase.On("Execute", mock.Anything, mock.Anything).Return(expectedSongs, nil)

	router := setupGetSongsRouter(mockLogger, mockUseCase)

	req, _ := http.NewRequest(http.MethodGet, "/songs?band=Test%20Group&limit=10", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response []SongResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, expectedSongs[0].Band, response[0].Band)
	assert.Equal(t, "2022-02-01", response[0].ReleaseDate)

	mockUseCase.AssertExpectations(t)
}

// Проверка на ошибку 400 при неверных параметрах запроса
func TestSongsHandler_GetSongsList_InvalidParams(t *testing.T) {
	mockLogger := new(MockLogger)
	mockUseCase := new(MockGetSongListUseCase)

	mockLogger.On("Debug", "Failed parsing request params", mock.Anything).Once()

	router := setupGetSongsRouter(mockLogger, mockUseCase)

	req, _ := http.NewRequest(http.MethodGet, "/songs?limit=invalid", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	mockLogger.AssertExpectations(t)
}

// Проверка на ошибку 404 когда песни не найдены
func TestSongsHandler_GetSongsList_NotFound(t *testing.T) {
	mockLogger := new(MockLogger)
	mockUseCase := new(MockGetSongListUseCase)

	mockLogger.On("Debug", "No songs found", mock.Anything).Once()
	mockUseCase.On("Execute", mock.Anything, mock.Anything).Return(nil, errs.ErrNotFound)

	router := setupGetSongsRouter(mockLogger, mockUseCase)

	req, _ := http.NewRequest(http.MethodGet, "/songs?song=Nonexistent%20Song", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusNotFound, recorder.Code)

	mockLogger.AssertExpectations(t)
	mockUseCase.AssertExpectations(t)
}

// Проверка на внутреннюю ошибку сервера
func TestSongsHandler_GetSongsList_InternalError(t *testing.T) {
	mockLogger := new(MockLogger)
	mockUseCase := new(MockGetSongListUseCase)

	mockLogger.On("Error", "Getting song list failed", mock.Anything).Once()
	mockUseCase.On("Execute", mock.Anything, mock.Anything).Return(nil, errors.New("database error"))

	router := setupGetSongsRouter(mockLogger, mockUseCase)

	req, _ := http.NewRequest(http.MethodGet, "/songs", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	mockLogger.AssertExpectations(t)
	mockUseCase.AssertExpectations(t)
}

// Проверка фильтрации по дате выпуска
func TestSongsHandler_GetSongsList_DateFilters(t *testing.T) {
	mockLogger := new(MockLogger)
	mockUseCase := new(MockGetSongListUseCase)

	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()

	releaseDate := time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC)
	expectedSongs := []entities.SongData{
		{
			ID:          123,
			Band:        "Test Group",
			Song:        "Test Song",
			ReleaseDate: releaseDate,
			Link:        "https://example.com/song1",
		},
	}

	mockUseCase.On("Execute", mock.Anything, mock.Anything).Return(expectedSongs, nil)

	router := setupGetSongsRouter(mockLogger, mockUseCase)

	req, _ := http.NewRequest(http.MethodGet, "/songs?release_date_from=2021-01-01&release_date_to=2023-01-01", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response []SongResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, "2022-02-01", response[0].ReleaseDate)

	mockUseCase.AssertExpectations(t)
}

// Проверка с использованием pagination (offset и limit)
func TestSongsHandler_GetSongsList_Pagination(t *testing.T) {
	mockLogger := new(MockLogger)
	mockUseCase := new(MockGetSongListUseCase)

	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()

	releaseDate := time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC)
	expectedSongs := []entities.SongData{
		{
			ID:          123,
			Band:        "Test Group",
			Song:        "Test Song",
			ReleaseDate: releaseDate,
			Link:        "https://example.com/song1",
		},
	}

	mockUseCase.On("Execute", mock.Anything, mock.Anything).Return(expectedSongs, nil)

	router := setupGetSongsRouter(mockLogger, mockUseCase)

	req, _ := http.NewRequest(http.MethodGet, "/songs?offset=10&limit=5", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response []SongResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, "2022-02-01", response[0].ReleaseDate)

	mockUseCase.AssertExpectations(t)
}
