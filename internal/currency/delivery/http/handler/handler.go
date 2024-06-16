package handler

import (
	"github.com/alemax1/currencies-api/internal/currency/service"
	"github.com/alemax1/currencies-api/pkg/logger"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	*service.Service
	Logger logger.Logger
}

func New(
	service *service.Service,
	l logger.Logger,
) *Handler {
	return &Handler{
		Service: service,
		Logger:  l,
	}
}

func (h Handler) InitRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	currencyApi := api.Group("/currency")

	currencyApi.Post("", h.CreateCurrency)
	currencyApi.Get("/rate", h.GetRate)
	currencyApi.Patch("/availability", h.ChangeCurrencyAvailability)
	currencyApi.Get("/all", h.GeteCurrencies)
}
