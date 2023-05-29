package service

import (
	"context"
	"fmt"

	"github.com/RipperAcskt/innotaxianalyst/config"
	"github.com/RipperAcskt/innotaxianalyst/internal/client"
	"github.com/google/uuid"
)

type Repo interface {
}

type GRPCService interface {
	GetOrdersQuantity(ctx context.Context, analys client.AnalysType) (int, error)
	GetJWT(ctx context.Context, id uuid.UUID) (*client.Token, error)
}
type Service struct {
	repo   Repo
	client GRPCService
	cfg    *config.Config
}

func New(repo Repo, client GRPCService, cfg *config.Config) *Service {
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
