package client

import (
	"context"
	"fmt"

	"github.com/RipperAcskt/innotaxi/pkg/proto"
	"github.com/RipperAcskt/innotaxianalyst/config"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientUser struct {
	clientUser proto.UserServiceClient
	connUser   *grpc.ClientConn

	cfg *config.Config
}

type Token struct {
	AccessToken  string `json:"Access_Token"`
	RefreshToken string `json:"Refresh_Token"`
}

func NewClientUser(cfg *config.Config) (*ClientUser, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	connUser, err := grpc.Dial(cfg.GRPC_USER_SERVICE_HOST, opts...)
	if err != nil {
		return nil, fmt.Errorf("dial user failed: %w", err)
	}
	clientUser := proto.NewUserServiceClient(connUser)

	return &ClientUser{
		clientUser: clientUser,
		connUser:   connUser,
		cfg:        cfg}, nil
}

func (u *ClientUser) GetJWT(ctx context.Context, id uuid.UUID) (*Token, error) {
	request := &proto.Params{
		Type: "analyst",
	}
	response, err := u.clientUser.GetJWT(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("get jwt failed: %w", err)
	}
	return &Token{response.AccessToken, response.RefreshToken}, nil
}

func (u *ClientUser) Close() error {
	return u.connUser.Close()
}
