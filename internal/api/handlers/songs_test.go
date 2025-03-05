package handlers_test

import (
	"bytes"
	"em-library/internal/api/handlers"
	"em-library/internal/entities"
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

func setupRouter(mockLogger *MockLogger, mockUseCases *MockCreateSongUseCase) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	useCases := usecase.UseCases{
		CreateSong: mockUseCases,
	}

	handler := handlers.NewSongsHandler(mockLogger, useCases)
	r.POST("/songs", handler.CreateSong)
	return r
}

// Удалось создать песню
func TestSongsHandler_CreateSong_Success(t *testing.T) {
	mockLogger := new(MockLogger)
	mockUseCase := new(MockCreateSongUseCase)

	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Error", mock.Anything, mock.Anything).Maybe()

	inputData := handlers.CreateSongParams{
		Group: "Test Group",
		Song:  "Test Song",
	}

	releaseDate := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	expectedSong := &entities.SongData{
		ID:          123,
		Group:       "Test Group",
		Song:        "Test Song",
		ReleaseDate: releaseDate,
		Link:        "https://example.com/song",
	}

	mockUseCase.On("Execute", mock.Anything, entities.NewSongData{
		Group: inputData.Group,
		Song:  inputData.Song,
	}).Return(expectedSong, nil)

	router := setupRouter(mockLogger, mockUseCase)

	jsonData, _ := json.Marshal(inputData)
	req, _ := http.NewRequest(http.MethodPost, "/songs", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusCreated, recorder.Code)

	var response handlers.SongResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, expectedSong.ID, response.ID)
	assert.Equal(t, expectedSong.Group, response.Group)
	assert.Equal(t, expectedSong.Song, response.Song)
	assert.Equal(t, "01.01.2022", response.ReleasedDate)
	assert.Equal(t, expectedSong.Link, response.Link)

	mockUseCase.AssertExpectations(t)
}

// Должен отдать BadRequest если передали битый JSON
func TestSongsHandler_CreateSong_InvalidJSON(t *testing.T) {
	mockLogger := new(MockLogger)
	mockUseCase := new(MockCreateSongUseCase)

	mockLogger.On("Debug", "Failed parsing credit request", mock.Anything).Once()

	router := setupRouter(mockLogger, mockUseCase)

	req, _ := http.NewRequest(http.MethodPost, "/songs", bytes.NewBuffer([]byte(`{invalid json}`)))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	mockLogger.AssertExpectations(t)
	mockUseCase.AssertNotCalled(t, "Execute")
}

// Должен отдать BadRequest если не хватает данных в полях
func TestSongsHandler_CreateSong_MissingFields(t *testing.T) {
	testCases := []struct {
		name        string
		requestBody map[string]string
		description string
	}{
		{
			name:        "Missing group field",
			requestBody: map[string]string{"song": "Test Song"},
			description: "Group field is completely missing",
		},
		{
			name:        "Empty group field",
			requestBody: map[string]string{"group": "", "song": "Test Song"},
			description: "Group field is empty string",
		},
		{
			name:        "Missing song field",
			requestBody: map[string]string{"group": "Test Group"},
			description: "Song field is completely missing",
		},
		{
			name:        "Empty song field",
			requestBody: map[string]string{"group": "Test Group", "song": ""},
			description: "Song field is empty string",
		},
		{
			name:        "Both fields missing",
			requestBody: map[string]string{},
			description: "Both group and song fields are missing",
		},
		{
			name:        "Both fields empty",
			requestBody: map[string]string{"group": "", "song": ""},
			description: "Both group and song fields are empty strings",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockLogger := new(MockLogger)
			mockUseCase := new(MockCreateSongUseCase)

			mockLogger.On("Debug", "Failed parsing credit request", mock.Anything).Once()

			router := setupRouter(mockLogger, mockUseCase)

			jsonData, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/songs", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusBadRequest, recorder.Code, tc.description)

			mockLogger.AssertExpectations(t)
			mockUseCase.AssertNotCalled(t, "Execute")
		})
	}
}

// если юзкейс вернул ошибку, то проверяем что ручка отдала 500-ку
func TestSongsHandler_CreateSong_UseCaseError(t *testing.T) {
	mockLogger := new(MockLogger)
	mockUseCase := new(MockCreateSongUseCase)

	mockLogger.On("Error", "Creation of song failed", mock.Anything).Once()

	inputData := handlers.CreateSongParams{
		Group: "Test Group",
		Song:  "Test Song",
	}

	mockUseCase.On("Execute", mock.Anything, entities.NewSongData{
		Group: inputData.Group,
		Song:  inputData.Song,
	}).Return(nil, errors.New("usecase error"))

	router := setupRouter(mockLogger, mockUseCase)

	jsonData, _ := json.Marshal(inputData)
	req, _ := http.NewRequest(http.MethodPost, "/songs", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	mockLogger.AssertExpectations(t)
	mockUseCase.AssertExpectations(t)
}
