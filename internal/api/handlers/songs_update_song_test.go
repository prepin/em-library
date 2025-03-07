package handlers_test

import (
	"bytes"
	"em-library/internal/api/handlers"
	"em-library/internal/entities"
	"em-library/internal/errs"
	"em-library/internal/usecase"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupPatchSongRouter(mockLogger *MockLogger, mockUseCase *MockUpdateSongUseCase) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	useCases := usecase.UseCases{
		UpdateSong: mockUseCase,
	}

	handler := handlers.NewSongsHandler(mockLogger, useCases)
	r.PATCH("/songs/:id", handler.UpdateSong)
	return r
}

func TestSongsHandler_UpdateSong_Success(t *testing.T) {
	mockLogger := new(MockLogger)
	mockUseCase := new(MockUpdateSongUseCase)

	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()

	songID := 123
	releaseDate := "2022-02-01"
	parsedDate, _ := time.Parse("2006-01-02", releaseDate)

	inputData := map[string]any{
		"band":         "Updated Band",
		"song":         "Updated Song",
		"release_date": releaseDate,
		"link":         "https://example.com/updated",
		"lyrics":       "Updated lyrics",
	}

	bandPtr := "Updated Band"
	songPtr := "Updated Song"
	linkPtr := "https://example.com/updated"
	lyricsPtr := "Updated lyrics"

	mockUseCase.On("Execute", mock.Anything, songID, entities.UpdateSongData{
		Band:        &bandPtr,
		Song:        &songPtr,
		ReleaseDate: &parsedDate,
		Link:        &linkPtr,
		Lyrics:      &lyricsPtr,
	}).Return(nil)

	router := setupPatchSongRouter(mockLogger, mockUseCase)

	jsonData, _ := json.Marshal(inputData)
	req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/songs/%d", songID), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusNoContent, recorder.Code)
	mockUseCase.AssertExpectations(t)
}

func TestSongsHandler_UpdateSong_InvalidID(t *testing.T) {
	mockLogger := new(MockLogger)
	mockUseCase := new(MockUpdateSongUseCase)

	mockLogger.On("Debug", "Missing or invalid ID param for request", mock.Anything).Once()

	router := setupPatchSongRouter(mockLogger, mockUseCase)

	req, _ := http.NewRequest(http.MethodPatch, "/songs/invalid", bytes.NewBuffer([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	mockLogger.AssertExpectations(t)
}

func TestSongsHandler_UpdateSong_InvalidJSON(t *testing.T) {
	mockLogger := new(MockLogger)
	mockUseCase := new(MockUpdateSongUseCase)

	mockLogger.On("Debug", "Failed parsing request params", mock.Anything).Once()

	router := setupPatchSongRouter(mockLogger, mockUseCase)

	req, _ := http.NewRequest(http.MethodPatch, "/songs/123", bytes.NewBuffer([]byte(`{invalid json}`)))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	mockLogger.AssertExpectations(t)
}

func TestSongsHandler_UpdateSong_InvalidFields(t *testing.T) {
	testCases := []struct {
		name        string
		requestBody map[string]string
		description string
	}{
		{
			name:        "Empty band field",
			requestBody: map[string]string{"band": ""},
			description: "Band field is empty string",
		},
		{
			name:        "Empty song field",
			requestBody: map[string]string{"song": ""},
			description: "Song field is empty string",
		},
		{
			name:        "Empty link field",
			requestBody: map[string]string{"link": ""},
			description: "Link field is empty string",
		},
		{
			name:        "Empty lyrics field",
			requestBody: map[string]string{"lyrics": ""},
			description: "Lyrics field is empty string",
		},
		{
			name:        "Invalid release date format",
			requestBody: map[string]string{"release_date": "invalid-date"},
			description: "Release date has invalid format",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockLogger := new(MockLogger)
			mockUseCase := new(MockUpdateSongUseCase)

			mockLogger.On("Debug", "Failed parsing request params", mock.Anything).Once()

			router := setupPatchSongRouter(mockLogger, mockUseCase)

			jsonData, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest(http.MethodPatch, "/songs/123", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusBadRequest, recorder.Code, tc.description)
			mockLogger.AssertExpectations(t)
		})
	}
}

func TestSongsHandler_UpdateSong_NotFound(t *testing.T) {
	mockLogger := new(MockLogger)
	mockUseCase := new(MockUpdateSongUseCase)

	mockLogger.On("Debug", "Song not found", mock.Anything).Once()

	songID := 123
	inputData := handlers.PatchSongParams{
		Band: stringPtr("Updated Band"),
	}

	mockUseCase.On("Execute", mock.Anything, songID, mock.Anything).Return(errs.ErrNotFound)

	router := setupPatchSongRouter(mockLogger, mockUseCase)

	jsonData, _ := json.Marshal(inputData)
	req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/songs/%d", songID), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
	mockLogger.AssertExpectations(t)
	mockUseCase.AssertExpectations(t)
}

func TestSongsHandler_UpdateSong_ServerError(t *testing.T) {
	mockLogger := new(MockLogger)
	mockUseCase := new(MockUpdateSongUseCase)

	mockLogger.On("Error", "Failed to update song", mock.Anything, mock.Anything, mock.Anything).Once()

	songID := 123
	inputData := handlers.PatchSongParams{
		Band: stringPtr("Updated Band"),
	}

	mockUseCase.On("Execute", mock.Anything, songID, mock.Anything).Return(errors.New("database error"))

	router := setupPatchSongRouter(mockLogger, mockUseCase)

	jsonData, _ := json.Marshal(inputData)
	req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/songs/%d", songID), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	mockLogger.AssertExpectations(t)
	mockUseCase.AssertExpectations(t)
}

func TestSongsHandler_UpdateSong_PartialUpdate(t *testing.T) {
	mockLogger := new(MockLogger)
	mockUseCase := new(MockUpdateSongUseCase)

	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()

	songID := 123
	inputData := handlers.PatchSongParams{
		Band: stringPtr("Updated Band"),
	}

	mockUseCase.On("Execute", mock.Anything, songID, entities.UpdateSongData{
		Band:        inputData.Band,
		Song:        nil,
		ReleaseDate: nil,
		Link:        nil,
		Lyrics:      nil,
	}).Return(nil)

	router := setupPatchSongRouter(mockLogger, mockUseCase)

	jsonData, _ := json.Marshal(inputData)
	req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/songs/%d", songID), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusNoContent, recorder.Code)
	mockUseCase.AssertExpectations(t)
}

func stringPtr(s string) *string {
	return &s
}
