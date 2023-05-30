package service

import (
	"context"
	"fmt"

	"github.com/RipperAcskt/innotaxianalyst/config"
	"github.com/RipperAcskt/innotaxianalyst/internal/client"
	"github.com/google/uuid"
)

const (
	user   = "user"
	driver = "driver"
)

type Repo interface {
}

type GRPCService interface {
	GetOrdersQuantity(ctx context.Context, analys client.AnalysType) (int, error)
	GetJWT(ctx context.Context, id uuid.UUID) (*client.Token, error)
	GetUserRating(ctx context.Context) ([]client.Rating, error)
	GetDriverRating(ctx context.Context) ([]client.Rating, error)
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

func (s *Service) GetRating(ctx context.Context, ratingType string) ([]client.Rating, error) {
	switch ratingType {
	case user:
		return s.client.GetUserRating(ctx)
	case driver:
		return s.client.GetDriverRating(ctx)
	default:
		return nil, fmt.Errorf("unknown type")
	}
}
