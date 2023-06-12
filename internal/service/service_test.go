package service_test

import (
	"context"
	"testing"

	"github.com/RipperAcskt/innotaxi/pkg/proto"
	"github.com/RipperAcskt/innotaxianalyst/config"
	"github.com/RipperAcskt/innotaxianalyst/internal/client"
	"github.com/RipperAcskt/innotaxianalyst/internal/model"
	"github.com/RipperAcskt/innotaxianalyst/internal/service"
	"github.com/RipperAcskt/innotaxianalyst/internal/service/mocks"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
)

func TestGetOrderAmount(t *testing.T) {
	type mockBehavior func(s *mocks.MockGRPCService)
	test := []struct {
		name         string
		mockBehavior mockBehavior
		err          error
	}{
		{
			name: "get order amount",
			mockBehavior: func(s *mocks.MockGRPCService) {
				s.EXPECT().GetOrdersQuantity(context.Background(), client.AnalysType("day")).Return(5, nil)
			},
			err: nil,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			grpc := mocks.NewMockGRPCService(ctrl)

			tt.mockBehavior(grpc)

			service := service.Service{
				Client: grpc,
			}

			_, err := service.GetOrderAmount(context.Background(), client.AnalysType("day"))
			assert.Equal(t, err, tt.err)
		})
	}
}

func TestSetRating(t *testing.T) {
	type mockBehavior func(repo *mocks.MockRepo, grpc *mocks.MockGRPCService)

	test := []struct {
		name         string
		rating       model.Rating
		mockBehavior mockBehavior
		err          error
	}{
		{
			name: "set rating",
			rating: model.Rating{
				Type:   "driver",
				ID:     "123",
				Rating: 4.2,
			},
			mockBehavior: func(repo *mocks.MockRepo, grpc *mocks.MockGRPCService) {
				repo.EXPECT().SetRatingUser(context.Background(), model.Rating{
					Type:   "driver",
					ID:     "123",
					Rating: 4.2,
				}).Return(4.2, nil)

				grpc.EXPECT().SetRating(context.Background(), &proto.Rating{
					Type: model.DriverType.ToString(),
					ID:   "123",
					Mark: 4.2,
				}).Return(&proto.Empty{}, nil)
			},
			err: nil,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mocks.NewMockRepo(ctrl)
			grpc := mocks.NewMockGRPCService(ctrl)

			tt.mockBehavior(repo, grpc)

			service := service.Service{
				Repo:   repo,
				Client: grpc,
			}

			err := service.SetRating(context.Background(), tt.rating)
			assert.Equal(t, err, tt.err)
		})
	}
}

func TestVerify(t *testing.T) {
	cfg := &config.Config{
		HS256_SECRET: "QWERTfg53gxb2",
	}

	test := []struct {
		name   string
		token  string
		userId uint64
		err    error
	}{
		{
			name:  "verify token expired",
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzY4Nzk5NDIsInR5cGUiOiJ1c2VyIiwidXNlcl9pZCI6MX0.qwiL4bupjm9O-ZnKpIcB8-erQytBJgkWlxnwPmRmv-c",
			err:   service.ErrTokenExpired,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			err := service.Verify(tt.token, cfg)
			assert.IsEqual(err, tt.err)
		})
	}
}
