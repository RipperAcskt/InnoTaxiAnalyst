package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/RipperAcskt/innotaxi/pkg/proto"
	"github.com/RipperAcskt/innotaxianalyst/config"
	"github.com/RipperAcskt/innotaxianalyst/internal/model"
	"github.com/RipperAcskt/innotaxianalyst/internal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	listener   net.Listener
	grpcServer *grpc.Server
	log        *zap.Logger
	service    *service.Service
	cfg        *config.Config
}

func New(log *zap.Logger, s *service.Service, cfg *config.Config) *Server {
	return &Server{nil, nil, log, s, cfg}
}

func (s *Server) Run() error {
	listener, err := net.Listen("tcp", s.cfg.GRPC_HOST)

	if err != nil {
		return fmt.Errorf("listen failed: %w", err)
	}

	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)

	s.listener = listener
	s.grpcServer = grpcServer

	proto.RegisterAnalystServiceServer(grpcServer, s)
	err = grpcServer.Serve(listener)
	if err != nil {
		return fmt.Errorf("serve failed: %w", err)
	}

	return nil
}

func (s *Server) SetRating(c context.Context, params *proto.Rating) (*proto.Empty, error) {
	r := model.Rating{
		Type:   params.Type,
		ID:     params.ID,
		Rating: params.Mark,
	}
	err := s.service.SetRating(c, r)
	return &proto.Empty{}, err
}

func (s *Server) Stop() error {
	s.log.Info("Shuttig down grpc...")

	err := s.listener.Close()
	if err != nil {
		return fmt.Errorf("listener close failed: %w", err)
	}

	s.grpcServer.Stop()
	s.log.Info("Grpc server exiting.")
	return nil
}
