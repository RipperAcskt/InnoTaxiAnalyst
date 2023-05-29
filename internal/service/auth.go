package service

import (
	"context"
	"fmt"
	"time"

	"github.com/RipperAcskt/innotaxianalyst/config"
	"github.com/RipperAcskt/innotaxianalyst/internal/client"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

var (
	ErrIncorrectLoginOrPassword = fmt.Errorf("incorrect login or password")
	ErrTokenExpired             = fmt.Errorf("token expired")
	ErrUnknownType              = fmt.Errorf("unknown type")
)

type SingIn struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (s *Service) SingIn(ctx context.Context, info *SingIn) (*client.Token, error) {
	if info.Login != s.cfg.ADMIN_LOGIN || info.Password != s.cfg.ADMIN_PASS {
		return nil, ErrIncorrectLoginOrPassword
	}
	return s.client.GetJWT(ctx, uuid.Nil)
}

func Verify(token string, cfg *config.Config) error {
	tokenJwt, err := jwt.Parse(
		token,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.HS256_SECRET), nil
		},
	)

	if err != nil {
		return fmt.Errorf("token parse failed: %w", err)
	}

	claims, ok := tokenJwt.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("jwt map claims failed")
	}

	if !claims.VerifyExpiresAt(time.Now().UTC().Unix(), true) {
		return ErrTokenExpired
	}
	if string(claims["type"].(string)) != "analyst" {
		return ErrUnknownType
	}
	return nil
}

func (s *Service) Refresh(ctx context.Context) (*client.Token, error) {
	token, err := s.client.GetJWT(ctx, uuid.Nil)
	if err != nil {
		return nil, fmt.Errorf("get jwt failed: %w", err)
	}

	return token, nil
}
