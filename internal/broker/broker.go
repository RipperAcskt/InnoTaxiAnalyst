package broker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/RipperAcskt/innotaxianalyst/config"
	"github.com/RipperAcskt/innotaxianalyst/internal/model"
	"github.com/segmentio/kafka-go"
)

type Broker struct {
	userReaded   *kafka.Reader
	driverReaded *kafka.Reader
	orderReaded  *kafka.Reader

	InfoChan chan infoStruct
	ErrChan  chan error

	cfg *config.Config
}

type infoStruct struct {
	InfoType model.ModelType
	Body     interface{}
}

func New(cfg *config.Config) (*Broker, error) {
	userReaded := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{cfg.BROKER_HOST},
		Topic:    model.UserType.ToString(),
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	userReaded.SetOffset(kafka.LastOffset)

	driverReaded := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{cfg.BROKER_HOST},
		Topic:    model.DriverType.ToString(),
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	driverReaded.SetOffset(kafka.LastOffset)

	orderReaded := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{cfg.BROKER_HOST},
		Topic:    model.OrderType.ToString(),
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	orderReaded.SetOffset(kafka.LastOffset)

	b := Broker{
		userReaded:   userReaded,
		driverReaded: driverReaded,
		orderReaded:  orderReaded,

		ErrChan:  make(chan error, 1),
		InfoChan: make(chan infoStruct, 1),

		cfg: cfg,
	}

	b.Read()
	return &b, nil
}

func (b *Broker) Read() {
	go func() {
		for {
			message, err := b.userReaded.ReadMessage(context.Background())
			if err != nil {
				b.ErrChan <- fmt.Errorf("read user failed: %w", err)
				continue
			}

			var user model.User
			err = json.Unmarshal(message.Value, &user)
			if err != nil {
				b.ErrChan <- fmt.Errorf("unmarshal user failed: %w", err)
				continue
			}

			info := infoStruct{
				InfoType: model.UserType,
				Body:     user,
			}
			b.InfoChan <- info
		}
	}()

	go func() {
		for {
			message, err := b.driverReaded.ReadMessage(context.Background())
			if err != nil {
				b.ErrChan <- fmt.Errorf("read user failed: %w", err)
				continue
			}

			var driver model.Driver
			err = json.Unmarshal(message.Value, &driver)
			if err != nil {
				b.ErrChan <- fmt.Errorf("unmarshal user failed: %w", err)
				continue
			}

			info := infoStruct{
				InfoType: model.DriverType,
				Body:     driver,
			}
			b.InfoChan <- info
		}
	}()

	go func() {
		for {
			message, err := b.orderReaded.ReadMessage(context.Background())
			if err != nil {
				b.ErrChan <- fmt.Errorf("read user failed: %w", err)
				continue
			}

			var order model.Order
			err = json.Unmarshal(message.Value, &order)
			if err != nil {
				b.ErrChan <- fmt.Errorf("unmarshal user failed: %w", err)
				continue
			}

			info := infoStruct{
				InfoType: model.OrderType,
				Body:     order,
			}
			b.InfoChan <- info
		}
	}()
}
