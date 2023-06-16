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

type ClientOrder struct {
	clientOrder proto.OrderServiceClient
	connOrder   *grpc.ClientConn

	cfg *config.Config
}

func NewClientOrder(cfg *config.Config) (*ClientOrder, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	connOrder, err := grpc.Dial(cfg.GRPC_ORDER_SERVICE_HOST, opts...)
	if err != nil {
		return nil, fmt.Errorf("dial order failed: %w", err)
	}
	clientOrder := proto.NewOrderServiceClient(connOrder)

	return &ClientOrder{
		clientOrder: clientOrder,
		connOrder:   connOrder,

		cfg: cfg}, nil
}

func (u *ClientOrder) GetOrdersQuantity(ctx context.Context, analys AnalysType) (int, error) {
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
	response, err := u.clientOrder.GetOrderQuantity(ctx, request)
	if err != nil {
		return 0, fmt.Errorf("get order quantity failed: %w", err)
	}
	return int(response.NumberOfOrders), nil
}

func (u *ClientOrder) Close() error {
	return u.connOrder.Close()
}
