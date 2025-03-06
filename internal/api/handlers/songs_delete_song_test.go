package handlers_test

import (
	"em-library/internal/api/handlers"
	"em-library/internal/usecase"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupDeleteSongRouter(mockLogger *MockLogger, mockUseCases *MockDeleteSongUseCase) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	useCases := usecase.UseCases{
		DeleteSong: mockUseCases,
	}

	handler := handlers.NewSongsHandler(mockLogger, useCases)
	r.DELETE("/songs/:id", handler.DeleteSong)
	return r
}

func TestSongsHandler_DeleteSong_Success(t *testing.T) {
	mockLogger := new(MockLogger)
	mockUseCase := new(MockDeleteSongUseCase)

	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()
	mockUseCase.On("Execute", mock.Anything, 123).Return(nil)

	router := setupDeleteSongRouter(mockLogger, mockUseCase)

	req, _ := http.NewRequest(http.MethodDelete, "/songs/123", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusNoContent, recorder.Code)
	mockUseCase.AssertExpectations(t)
}

func TestSongsHandler_DeleteSong_InvalidID(t *testing.T) {
	mockLogger := new(MockLogger)
	mockUseCase := new(MockDeleteSongUseCase)

	mockLogger.On("Debug", "Missing or invalid ID param for request", mock.Anything).Once()

	router := setupDeleteSongRouter(mockLogger, mockUseCase)

	req, _ := http.NewRequest(http.MethodDelete, "/songs/invalid", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	mockLogger.AssertExpectations(t)
}

func TestSongsHandler_DeleteSong_InternalError(t *testing.T) {
	mockLogger := new(MockLogger)
	mockUseCase := new(MockDeleteSongUseCase)

	mockLogger.On("Debug", "Failed to delete song", mock.Anything).Once()
	mockUseCase.On("Execute", mock.Anything, 123).Return(assert.AnError)

	router := setupDeleteSongRouter(mockLogger, mockUseCase)

	req, _ := http.NewRequest(http.MethodDelete, "/songs/123", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	mockLogger.AssertExpectations(t)
	mockUseCase.AssertExpectations(t)
}
