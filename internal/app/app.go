package app

import (
	"fmt"
	"net/http"

	"github.com/RipperAcskt/innotaxianalyst/config"
	"github.com/RipperAcskt/innotaxianalyst/internal/client"
	"github.com/RipperAcskt/innotaxianalyst/internal/handler"
	"github.com/RipperAcskt/innotaxianalyst/internal/repo/clickhouse"
	"github.com/RipperAcskt/innotaxianalyst/internal/server"
	"github.com/RipperAcskt/innotaxianalyst/internal/service"
	"go.uber.org/zap"
)

func Run() error {
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("config new failed: %w", err)
	}

	log, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("new production failed: %w", err)
	}

	repo, err := clickhouse.New(cfg)
	if err != nil {
		return fmt.Errorf("clickhouse new failed: %w", err)
	}

	clientUser, err := client.NewClientUser(cfg)
	if err != nil {
		return fmt.Errorf("client user new failed: %w", err)
	}

	clientOrder, err := client.NewClientOrder(cfg)
	if err != nil {
		return fmt.Errorf("client user new failed: %w", err)
	}

	service := service.New(repo, clientUser, clientOrder, cfg)

	handler := handler.New(service, cfg, log)

	server := &server.Server{
		Log: log,
	}

	if err := server.Run(handler.InitRouters(), cfg); err != nil && err != http.ErrServerClosed {
		log.Error(fmt.Sprintf("server run failed: %v", err))
	}

	if err := server.ShutDown(); err != nil {
		return fmt.Errorf("server shut down failed: %w", err)
	}
	return nil
}
