package client

import (
	"context"
	"fmt"
	"time"

	"github.com/RipperAcskt/innotaxi/pkg/proto"
	"github.com/RipperAcskt/innotaxianalyst/config"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AnalysType string

const (
	day   AnalysType = "day"
	month AnalysType = "month"
)

type User struct {
	clientOrder proto.OrderServiceClient
	connOrder   *grpc.ClientConn
	clientUser  proto.UserServiceClient
	connUser    *grpc.ClientConn
	cfg         *config.Config
}

type Token struct {
	AccessToken  string `json:"Access_Token"`
	RefreshToken string `json:"Refresh_Token"`
}

func New(cfg *config.Config) (*User, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	connOrder, err := grpc.Dial(cfg.GRPC_ORDER_SERVICE_HOST, opts...)
	if err != nil {
		return nil, fmt.Errorf("dial order failed: %w", err)
	}
	clientOrder := proto.NewOrderServiceClient(connOrder)

	connUser, err := grpc.Dial(cfg.GRPC_USER_SERVICE_HOST, opts...)
	if err != nil {
		return nil, fmt.Errorf("dial user failed: %w", err)
	}
	clientUser := proto.NewUserServiceClient(connUser)

	return &User{
		clientOrder: clientOrder,
		connOrder:   connOrder,
		clientUser:  clientUser,
		connUser:    connUser,
		cfg:         cfg}, nil
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
	response, err := u.clientOrder.GetOrderQuantity(ctx, request)
	if err != nil {
		return 0, fmt.Errorf("get order quantity failed: %w", err)
	}
	return int(response.NumberOfOrders), nil
}

func (u *User) GetJWT(ctx context.Context, id uuid.UUID) (*Token, error) {
	request := &proto.Params{
		Type: "analyst",
	}
	response, err := u.clientUser.GetJWT(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("get jwt failed: %w", err)
	}
	return &Token{response.AccessToken, response.RefreshToken}, nil
}

func (u *User) Close() error {
	if err := u.connOrder.Close(); err != nil {
		return err
	}
	return u.connUser.Close()
}
