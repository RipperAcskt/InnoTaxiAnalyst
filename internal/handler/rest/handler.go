package handler

import (
	"fmt"
	"net/http"

	"github.com/RipperAcskt/innotaxianalyst/config"
	"github.com/RipperAcskt/innotaxianalyst/internal/client"
	"github.com/RipperAcskt/innotaxianalyst/internal/model"
	"github.com/RipperAcskt/innotaxianalyst/internal/service"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Handler struct {
	service *service.Service
	cfg     *config.Config
	log     *zap.Logger
}

type ErrResp struct {
	Err string `json:"error"`
}

func New(s *service.Service, cfg *config.Config, log *zap.Logger) *Handler {
	return &Handler{
		service: s,
		cfg:     cfg,
		log:     log,
	}
}

func (h *Handler) InitRouters() *fiber.App {
	router := fiber.New()

	analyst := router.Group("/analyst")
	analyst.Post("/sing-in", h.SingIn)
	analyst.Use(h.VerifyToken())
	analyst.Post("/amount", h.GetOrdersAmount)
	analyst.Post("/rating", h.GetRating)
	analyst.Post("/refresh", h.Refresh)

	return router
}

func (h *Handler) GetOrdersAmount(ctx *fiber.Ctx) error {
	var analysType model.AnalysType
	if err := ctx.BodyParser(&analysType); err != nil {
		h.log.Sugar().Errorf("body parser failed: %w", err)
		ctx = ctx.Status(http.StatusBadRequest)
		return ctx.SendString(fmt.Sprintf("body parse failed: %v", err))
	}

	num, err := h.service.GetOrderAmount(ctx.UserContext(), client.AnalysType(analysType.AnalysType))
	if err != nil {
		h.log.Sugar().Errorf("get order amount failed: %w", err)
		ctx = ctx.Status(http.StatusInternalServerError)
		return ctx.SendString(fmt.Sprintf("get order amount failed: %v", err))
	}

	var resp struct {
		Amount int `json:"amount"`
	}
	resp.Amount = num

	ctx = ctx.Status(http.StatusOK)
	return ctx.JSON(resp)
}

func (h *Handler) GetRating(ctx *fiber.Ctx) error {
	var ratingType struct {
		RatingType string `json:"type"`
	}
	if err := ctx.BodyParser(&ratingType); err != nil {
		h.log.Sugar().Errorf("body parser failed: %w", err)
		ctx = ctx.Status(http.StatusBadRequest)
		return ctx.SendString(fmt.Sprintf("body parse failed: %v", err))
	}

	ratings, err := h.service.GetRating(ctx.UserContext(), ratingType.RatingType)
	if err != nil {
		h.log.Sugar().Errorf("get rating failed: %w", err)
		ctx = ctx.Status(http.StatusInternalServerError)
		return ctx.SendString(fmt.Sprintf("get rating failed: %v", err))
	}

	var resp struct {
		Ratings []model.Rating `json:"ratings"`
	}
	resp.Ratings = ratings

	ctx = ctx.Status(http.StatusOK)
	return ctx.JSON(resp)
}
