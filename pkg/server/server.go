package server

import (
	"context"
	"em-library/config"
	"em-library/internal/app/api/handlers"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	httpServer *http.Server
	cfg        *config.Config
}

func New(cfg *config.Config, handlers *handlers.Handlers) *Server {
	if cfg.Server.ProductionMode {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	handlers.RegisterRoutes(router)

	return &Server{
		httpServer: &http.Server{
			Addr:    ":" + cfg.Server.Port,
			Handler: router,
		},
		cfg: cfg,
	}
}

func (s *Server) Run() error {
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.cfg.Logger.Error("failed to launch server", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	return s.Shutdown()
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	select {
	case <-ctx.Done():
		s.cfg.Logger.Warn("timeout of 5 seconds.")
	default:
		s.cfg.Logger.Info("Server shutdown done.")
	}

	return nil
}
