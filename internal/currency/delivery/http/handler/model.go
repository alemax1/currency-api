package handler

import (
	"encoding/json"

	"github.com/alemax1/currencies-api/internal/currency/domain"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func ParseAndValidateRequest[T interface{ Validate() error }](body []byte) (T, error) {
	var req T
	if err := json.Unmarshal(body, &req); err != nil {
		return req, errInvalidJSONBodyRequest
	}

	return req, req.Validate()
}

type currencyCreateRequest struct {
	Name string              `json:"name"`
	Type domain.CurrencyType `json:"type"`
}

func (r currencyCreateRequest) Validate() error {
	if err := validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required, validation.Length(2, 255)),
		validation.Field(&r.Type, validation.In(domain.Crypto, domain.Fiat)),
	); err != nil {
		return errInvalidInput
	}

	return nil
}

type changeCurrencyAvailabilityRequest struct {
	Name        string `json:"name"`
	IsAvailable bool   `json:"isAvailable"`
}

func (r changeCurrencyAvailabilityRequest) Validate() error {
	if err := validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required, validation.Length(2, 255)),
	); err != nil {
		return errInvalidInput
	}

	return nil
}

type createCurrencyResponse struct {
	ID int64 `json:"id"`
}

type defaultResponse struct {
	Success bool `json:"success"`
}

type getRateResponse struct {
	Rate float64 `json:"rate"`
}

type errResponse struct {
	Error string `json:"err"`
}

type currency struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	ValueUSD    float64 `json:"valueUSD"`
	IsAvailable bool    `json:"isAvailable"`
}

type getAvailableCurrenciesResponse struct {
	Currencies []currency `json:"currencies"`
}

func currencyToDto(curr domain.Currency) currency {
	return currency{
		ID:          curr.ID,
		Name:        curr.Name,
		Type:        string(curr.Type),
		ValueUSD:    curr.ValueUSD.InexactFloat64(),
		IsAvailable: curr.IsAvailable,
	}
}
