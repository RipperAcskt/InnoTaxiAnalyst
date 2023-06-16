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

type UserService interface {
	GetJWT(ctx context.Context, id uuid.UUID) (*client.Token, error)
}

type OrderService interface {
	GetOrdersQuantity(ctx context.Context, analys client.AnalysType) (int, error)
}
type Service struct {
	repo        Repo
	clientUser  UserService
	clientOrder OrderService
	cfg         *config.Config
}

func New(repo Repo, clientUser UserService, clientOrder OrderService, cfg *config.Config) *Service {
	return &Service{
		repo:        repo,
		clientUser:  clientUser,
		clientOrder: clientOrder,
		cfg:         cfg,
	}
}

func (s *Service) GetOrderAmount(ctx context.Context, analys client.AnalysType) (int, error) {
	num, err := s.clientOrder.GetOrdersQuantity(ctx, analys)
	if err != nil {
		return 0, fmt.Errorf("get orders quantity failed: %w", err)
	}
	return num, nil
}
