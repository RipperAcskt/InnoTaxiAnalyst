package client

import (
	"context"
	"fmt"

	"github.com/RipperAcskt/innotaxi/pkg/proto"
	"github.com/RipperAcskt/innotaxianalyst/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientDriver struct {
	clientDriver proto.DriverServiceClient
	connDriver   *grpc.ClientConn

	cfg *config.Config
}

func NewClientDriver(cfg *config.Config) (*ClientDriver, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	connDriver, err := grpc.Dial(cfg.GRPC_DRIVER_SERVICE_HOST, opts...)
	if err != nil {
		return nil, fmt.Errorf("dial user failed: %w", err)
	}
	clientDriver := proto.NewDriverServiceClient(connDriver)

	return &ClientDriver{
		clientDriver: clientDriver,
		connDriver:   connDriver,

		cfg: cfg}, nil
}

func (u *ClientDriver) SetRating(ctx context.Context, rating *proto.Rating) (*proto.Empty, error) {
	_, err := u.clientDriver.SetRating(ctx, rating)
	if err != nil {
		return nil, fmt.Errorf("set rating driver failed: %w", err)
	}

	return &proto.Empty{}, nil
}

func (u *ClientDriver) Close() error {
	return u.connDriver.Close()
}
