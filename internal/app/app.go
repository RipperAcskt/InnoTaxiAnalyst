package app

import (
	"fmt"
	"net/http"

	"github.com/RipperAcskt/innotaxianalyst/config"
	"github.com/RipperAcskt/innotaxianalyst/internal/broker"
	"github.com/RipperAcskt/innotaxianalyst/internal/client"
	"github.com/RipperAcskt/innotaxianalyst/internal/handler/grpc"
	handler "github.com/RipperAcskt/innotaxianalyst/internal/handler/rest"
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

	client, err := client.New(cfg)
	if err != nil {
		return fmt.Errorf("client new failed: %w", err)
	}

	broker, err := broker.New(cfg)
	if err != nil {
		return fmt.Errorf("broker new failed: %w", err)
	}

	go func() {
		err := <-broker.ErrChan
		log.Error("error: ", zap.Error(err))
	}()

	service := service.New(repo, client, broker, cfg)

	handler := handler.New(service, cfg, log)

	server := &server.Server{
		Log: log,
	}

	go func() {
		if err := server.Run(handler.InitRouters(), cfg); err != nil && err != http.ErrServerClosed {
			log.Error(fmt.Sprintf("server run failed: %v", err))
			return
		}
	}()

	grpcServer := grpc.New(log, service, cfg)
	go func() {
		if err := grpcServer.Run(); err != nil {
			log.Error(fmt.Sprintf("grpc server run failed: %v", err))
			return
		}
	}()

	if err := server.ShutDown(); err != nil {
		return fmt.Errorf("server shut down failed: %w", err)
	}
	return nil
}
