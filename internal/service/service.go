package service

import (
	"context"
	"fmt"

	"github.com/RipperAcskt/innotaxianalyst/config"
	"github.com/RipperAcskt/innotaxianalyst/internal/broker"
	"github.com/RipperAcskt/innotaxianalyst/internal/client"
	"github.com/RipperAcskt/innotaxianalyst/internal/model"
	"github.com/google/uuid"
)

type Repo interface {
	WriteUser(user model.User) error
	WriteDriver(driver model.Driver) error
	WriteOrder(order model.Order) error
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
	broker *broker.Broker
	cfg    *config.Config
}

func New(repo Repo, client GRPCService, broker *broker.Broker, cfg *config.Config) *Service {
	s := Service{
		repo:   repo,
		client: client,
		broker: broker,
		cfg:    cfg,
	}
	go s.GetMessages()
	return &s
}

func (s *Service) GetOrderAmount(ctx context.Context, analys client.AnalysType) (int, error) {
	num, err := s.client.GetOrdersQuantity(ctx, analys)
	if err != nil {
		return 0, fmt.Errorf("get orders quantity failed: %w", err)
	}
	return num, nil
}

func (s *Service) GetRating(ctx context.Context, ratingType string) ([]client.Rating, error) {
	modelType := model.New(ratingType)
	switch modelType {
	case model.UserType:
		return s.client.GetUserRating(ctx)
	case model.DriverType:
		return s.client.GetDriverRating(ctx)
	default:
		return nil, fmt.Errorf("unknown type")
	}
}

func (s *Service) GetMessages() {
	for {
		info := <-s.broker.InfoChan

		switch info.InfoType {
		case model.UserType:
			user := info.Body.(model.User)

			uuid, err := uuid.NewRandom()
			if err != nil {
				s.broker.ErrChan <- fmt.Errorf("uuid user new failed: %w", err)
			}
			user.ID = uuid

			err = s.repo.WriteUser(user)
			if err != nil {
				fmt.Println(err)
				s.broker.ErrChan <- fmt.Errorf("write user failed: %w", err)
			}

		case model.DriverType:
			driver := info.Body.(model.Driver)

			uuid, err := uuid.NewRandom()
			if err != nil {
				s.broker.ErrChan <- fmt.Errorf("uuid driver new failed: %w", err)
			}
			driver.ID = uuid

			err = s.repo.WriteDriver(driver)
			if err != nil {
				fmt.Println(err)
				s.broker.ErrChan <- fmt.Errorf("write driver failed: %w", err)
			}

		case model.OrderType:
			order := info.Body.(model.Order)

			uuid, err := uuid.NewRandom()
			if err != nil {
				s.broker.ErrChan <- fmt.Errorf("uuid order new failed: %w", err)
			}
			order.ID = uuid

			err = s.repo.WriteOrder(order)
			if err != nil {
				fmt.Println(err)
				s.broker.ErrChan <- fmt.Errorf("write order failed: %w", err)
			}
		}
	}
}
