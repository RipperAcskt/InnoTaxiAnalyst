package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/RipperAcskt/innotaxianalyst/internal/client"
	"github.com/RipperAcskt/innotaxianalyst/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func (h *Handler) SingIn(ctx *fiber.Ctx) error {
	var authInfo service.SingIn
	if err := ctx.BodyParser(&authInfo); err != nil {
		h.log.Sugar().Errorf("body parser failed: %w", err)
		ctx = ctx.Status(http.StatusBadRequest)
		return ctx.SendString(fmt.Sprintf("body parse failed: %v", err))
	}

	token, err := h.service.SingIn(ctx.Context(), &authInfo)
	if err != nil {
		h.log.Sugar().Errorf("sing in failed: %w", err)
		ctx = ctx.Status(http.StatusInternalServerError)
		return ctx.SendString(fmt.Sprintf("sing in failed: %v", err))
	}

	ctx = ctx.Status(http.StatusOK)
	return ctx.JSON(token)
}

func (h *Handler) VerifyToken(args ...interface{}) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		headers := ctx.GetReqHeaders()
		token := strings.Split(headers["Authorization"], " ")
		if len(token) < 2 {
			err := ErrResp{
				Err: fmt.Errorf("access token required").Error(),
			}
			ctx = ctx.Status(http.StatusBadRequest)
			return ctx.JSON(err)
		}
		accessToken := token[1]

		err := service.Verify(accessToken, h.cfg)
		if err != nil {
			if errors.Is(err, service.ErrTokenExpired) {
				err := ErrResp{
					Err: err.Error(),
				}
				ctx = ctx.Status(http.StatusUnauthorized)
				return ctx.JSON(err)
			}
			if strings.Contains(err.Error(), jwt.ErrSignatureInvalid.Error()) {
				err := ErrResp{
					Err: fmt.Errorf("wrong signature").Error(),
				}
				ctx = ctx.Status(http.StatusForbidden)
				return ctx.JSON(err)
			}

			err := ErrResp{
				Err: fmt.Errorf("verify failed: %w", err).Error(),
			}
			ctx = ctx.Status(http.StatusInternalServerError)
			return ctx.JSON(err)
		}

		return ctx.Next()
	}
}

func (h *Handler) Refresh(ctx *fiber.Ctx) error {
	var token client.Token
	if err := ctx.BodyParser(&token); err != nil {
		h.log.Sugar().Errorf("body parser failed: %w", err)
		ctx = ctx.Status(http.StatusBadRequest)
		return ctx.SendString(fmt.Sprintf("body parse failed: %v", err))
	}

	err := service.Verify(token.RefreshToken, h.cfg)
	if err != nil {
		if err != nil {
			if errors.Is(err, service.ErrTokenExpired) {
				err := ErrResp{
					Err: err.Error(),
				}
				ctx = ctx.Status(http.StatusUnauthorized)
				return ctx.JSON(err)
			}
			if strings.Contains(err.Error(), jwt.ErrSignatureInvalid.Error()) {
				err := ErrResp{
					Err: fmt.Errorf("wrong signature").Error(),
				}
				ctx = ctx.Status(http.StatusForbidden)
				return ctx.JSON(err)
			}

			err := ErrResp{
				Err: fmt.Errorf("verify failed: %w", err).Error(),
			}
			ctx = ctx.Status(http.StatusInternalServerError)
			return ctx.JSON(err)
		}
	}

	t, err := h.service.Refresh(ctx.Context())
	if err != nil {
		err := ErrResp{
			Err: fmt.Errorf("refresh failed: %w", err).Error(),
		}
		ctx = ctx.Status(http.StatusInternalServerError)
		return ctx.JSON(err)
	}

	ctx = ctx.Status(http.StatusOK)
	return ctx.JSON(t)
}
