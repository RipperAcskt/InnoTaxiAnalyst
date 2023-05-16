package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RipperAcskt/innotaxianalyst/config"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Server struct {
	app *fiber.App
	Log *zap.Logger
}

func (s *Server) Run(app *fiber.App, cfg *config.Config) error {
	s.app = app

	s.Log.Sugar().Infof("connect to http://%s/", cfg.SERVER_HOST)
	return app.Listen(cfg.SERVER_HOST)
}

func (s *Server) ShutDown() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	s.Log.Info("Shuttig down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.app.ShutdownWithContext(ctx); err != nil {
		return fmt.Errorf("shut down failed: %w", err)
	}

	s.Log.Info("Server exiting.")
	return nil
}
