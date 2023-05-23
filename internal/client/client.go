package client

import (
	"context"
	"fmt"
	"time"

	"github.com/RipperAcskt/innotaxi/pkg/proto"
	"github.com/RipperAcskt/innotaxianalyst/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AnalysType string

const (
	day   AnalysType = "day"
	month AnalysType = "month"
)

type User struct {
	client proto.OrderServiceClient
	conn   *grpc.ClientConn
	cfg    *config.Config
}

func New(cfg *config.Config) (*User, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.Dial(cfg.GRPC_ORDER_SERVICE_HOST, opts...)

	if err != nil {
		return nil, fmt.Errorf("dial failed: %w", err)
	}

	client := proto.NewOrderServiceClient(conn)

	return &User{client, conn, cfg}, nil
}

func (u *User) GetOrdersQuantity(ctx context.Context, analys AnalysType) (int, error) {
	var timeStr string
	timeNow := time.Now()
	if analys == day {
		timeDay := timeNow.AddDate(0, 0, -1)
		timeStr = timeDay.Format("2006-01-02 15:04:05")
	} else if analys == month {
		timeMonth := timeNow.AddDate(0, -1, 0)
		timeStr = timeMonth.Format("2006-01-02 15:04:05")
	}

	request := &proto.Time{
		TimeStarted: timeStr,
	}
	response, err := u.client.GetOrderQuantity(ctx, request)
	if err != nil {
		return 0, fmt.Errorf("get order quantity failed: %w", err)
	}
	return int(response.NumberOfOrders), nil
}

func (u *User) Close() error {
	return u.conn.Close()
}
