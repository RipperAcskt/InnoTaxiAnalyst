package service

import (
	"context"
	"fmt"

	"github.com/RipperAcskt/innotaxianalyst/config"
	"github.com/RipperAcskt/innotaxianalyst/internal/client"
)

type Repo interface {
}

type OrderService interface {
	GetOrdersQuantity(ctx context.Context, analys client.AnalysType) (int, error)
}

type Service struct {
	repo   Repo
	client OrderService
	cfg    *config.Config
}

func New(repo Repo, client OrderService, cfg *config.Config) *Service {
	return &Service{
		repo:   repo,
		client: client,
		cfg:    cfg,
	}
}

func (s *Service) GetOrderAmount(ctx context.Context, analys client.AnalysType) (int, error) {
	num, err := s.client.GetOrdersQuantity(ctx, analys)
	if err != nil {
		return 0, fmt.Errorf("get orders quantity failed: %w", err)
	}
	return num, nil
}
