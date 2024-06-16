package service

import (
	"context"
	"time"

	"github.com/alemax1/currencies-api/pkg/logger"
)

type Service struct {
	Currency *currency
}

func New(
	CurrencyRepo CurrencyRepo,
	ForexAPI ForexAPI,
	logger logger.Logger,
) *Service {
	currencySvc := newCurrency(
		CurrencyRepo,
		ForexAPI,
		logger,
	)

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := currencySvc.UpdateCryptoCurrencies(ctx); err != nil {
			currencySvc.Logger.Error().Err(err).Msgf("update crypto currencies")
		}
		if err := currencySvc.UpdateFiatCurrencies(ctx); err != nil {
			currencySvc.Logger.Error().Err(err).Msgf("update fiat currencies")
		}
	}()

	return &Service{
		Currency: currencySvc,
	}
}
