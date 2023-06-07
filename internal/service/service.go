package service

import (
	"context"
	"fmt"

	"github.com/RipperAcskt/innotaxi/pkg/proto"
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
	SetRatingUser(ctx context.Context, r model.Rating) (float64, error)
	SetRatingDriver(ctx context.Context, r model.Rating) (float64, error)
	GetRating(ctx context.Context, db string) ([]model.Rating, error)
}

type GRPCService interface {
	GetOrdersQuantity(ctx context.Context, analys client.AnalysType) (int, error)
	GetJWT(ctx context.Context, id uuid.UUID) (*client.Token, error)
	SetRating(c context.Context, params *proto.Rating) (*proto.Empty, error)
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

func (s *Service) SetRating(ctx context.Context, r model.Rating) error {
	rating := &proto.Rating{
		Type: r.Type,
		ID:   r.ID,
	}

	switch r.Type {
	case model.DriverType.ToString():
		rate, err := s.repo.SetRatingUser(ctx, r)
		if err != nil {
			return fmt.Errorf("set rating user failed: %w", err)
		}

		rating.Mark = float32(rate)

	case model.UserType.ToString():
		rate, err := s.repo.SetRatingDriver(ctx, r)
		if err != nil {
			return fmt.Errorf("set rating driver failed: %w", err)
		}

		rating.Mark = float32(rate)
	}

	_, err := s.client.SetRating(ctx, rating)
	if err != nil {
		return fmt.Errorf("set rating failed: %w", err)
	}
	return nil
}

func (s *Service) GetRating(ctx context.Context, ratingType string) ([]model.Rating, error) {
	return s.repo.GetRating(ctx, ratingType)
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
