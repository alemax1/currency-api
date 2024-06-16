package service

import (
	"context"
	"fmt"

	"github.com/alemax1/currencies-api/internal/currency/domain"
	"github.com/alemax1/currencies-api/pkg/logger"
	"github.com/shopspring/decimal"
)

type ForexAPI interface {
	SendFetchOneRequest(ctx context.Context, currency domain.Currency) (domain.CurrencyWithValue, error)
	SendMultiFetchRequest(ctx context.Context, currencies []domain.Currency) ([]domain.CurrencyWithValue, error)
}

type CurrencyRepo interface {
	AddEmptyCurrency(ctx context.Context, currency domain.Currency) (int64, error)
	GetCurrency(ctx context.Context, name string) (domain.Currency, error)
	UpdateCurrencyAvailability(ctx context.Context, name string, isAvailable bool) error
	GetAll(ctx context.Context) ([]domain.Currency, error)
	GetCurrenciesByType(ctx context.Context, tp domain.CurrencyType) ([]domain.Currency, error)
	UpdateCurrencyByName(ctx context.Context, currency domain.CurrencyUpdateData) error
}

type currency struct {
	CurrencyRepo CurrencyRepo
	ForexAPI     ForexAPI
	Logger       logger.Logger
}

func newCurrency(
	currencyRepo CurrencyRepo,
	forexAPI ForexAPI,
	logger logger.Logger,
) *currency {
	return &currency{
		CurrencyRepo: currencyRepo,
		ForexAPI:     forexAPI,
		Logger:       logger,
	}
}

func (c currency) Create(ctx context.Context, currency domain.Currency) (int64, error) {
	id, err := c.CurrencyRepo.AddEmptyCurrency(ctx, currency)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (c currency) GetRate(ctx context.Context, rate domain.Rate) (float64, error) {
	currencyFrom, err := c.CurrencyRepo.GetCurrency(ctx, rate.From)
	if err != nil {
		return 0, fmt.Errorf("get from currency: %w", err)
	}

	currencyTo, err := c.CurrencyRepo.GetCurrency(ctx, rate.To)
	if err != nil {
		return 0, fmt.Errorf("get to currency: %w", err)
	}

	if currencyFrom.Type == currencyTo.Type {
		return 0, domain.NewServiceError(domain.ErrInvalidCurrencyTypes, domain.Client)
	}

	if !currencyFrom.IsAvailable || !currencyTo.IsAvailable {
		return 0, domain.NewServiceError(domain.ErrInvalidCurrencyTypes, domain.Client)
	}

	if decimal.Zero.Equal(currencyFrom.ValueUSD) || decimal.Zero.Equal(currencyTo.ValueUSD) {
		return 0, domain.NewServiceError(domain.ErrValueCannotBeZero, domain.Client)
	}

	// The rate value of the currency is divided by its dollar equivalent and multiplied by the dollar equivalent of the currency to exchange
	rateValue := rate.Value.Div(currencyFrom.ValueUSD).Mul(currencyTo.ValueUSD)

	return rateValue.InexactFloat64(), nil
}

func (c currency) ChangeAvailability(ctx context.Context, name string, isAvailable bool) error {
	if err := c.CurrencyRepo.UpdateCurrencyAvailability(ctx, name, isAvailable); err != nil {
		return err
	}

	return nil
}

func (c currency) GetAll(ctx context.Context) ([]domain.Currency, error) {
	currencies, err := c.CurrencyRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return currencies, nil
}

func (c currency) UpdateFiatCurrencies(ctx context.Context) error {
	currencies, err := c.CurrencyRepo.GetCurrenciesByType(ctx, domain.Fiat)
	if err != nil {
		return fmt.Errorf("get all currencies by type: %w", err)
	}

	resp, err := c.ForexAPI.SendMultiFetchRequest(ctx, currencies)
	if err != nil || len(resp) == 0 {
		return fmt.Errorf("send multi fetch request: %w", err)
	}

	for _, currency := range resp {
		currencyUpdate := domain.CurrencyUpdateData{
			Name:        currency.Name,
			ValueUSD:    currency.Value,
			IsAvailable: true,
		}

		if err := c.CurrencyRepo.UpdateCurrencyByName(ctx, currencyUpdate); err != nil {
			c.Logger.Error().Err(err).Msgf("update currency, value:%s, name:%s", currency.Value, currency.Name)

			if err := c.CurrencyRepo.UpdateCurrencyAvailability(ctx, currency.Name, false); err != nil {
				c.Logger.Error().Err(err).Msgf("update currency availability name:%s", currency.Name)
			}
		}
	}

	return nil
}

func (c currency) UpdateCryptoCurrencies(ctx context.Context) error {
	currencies, err := c.CurrencyRepo.GetCurrenciesByType(context.Background(), domain.Crypto)
	if err != nil {
		return fmt.Errorf("get currencies by type: %w", err)
	}

	for _, currency := range currencies {
		resp, err := c.ForexAPI.SendFetchOneRequest(ctx, currency)
		if err != nil {
			c.Logger.Error().Err(err).Msg("send fetch one request")
			continue
		}

		currencyUpdate := domain.CurrencyUpdateData{
			Name:        resp.Name,
			ValueUSD:    resp.Value,
			IsAvailable: true,
		}

		if err := c.CurrencyRepo.UpdateCurrencyByName(context.Background(), currencyUpdate); err != nil {
			c.Logger.Error().Err(err).Msgf("update currency, value:%s, name:%s", resp.Value, currency.Name)

			if err := c.CurrencyRepo.UpdateCurrencyAvailability(ctx, currency.Name, false); err != nil {
				c.Logger.Error().Err(err).Msgf("update currency availability name:%s", currency.Name)
			}
		}
	}

	return nil
}
