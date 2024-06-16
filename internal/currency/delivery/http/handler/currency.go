package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/alemax1/currencies-api/internal/currency/domain"
	"github.com/gofiber/fiber/v3"
	"github.com/shopspring/decimal"
)

const (
	fromQueryParam  = "from"
	toQueryParam    = "to"
	valueQueryParam = "value"
)

// CreateCurrency godoc
//
//	@Summary		add a new currency
//	@Description	create currency
//	@Tags			currency
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	createCurrencyResponse
//	@Failure		400	{object}	errResponse
//	@Failure		500	{object}	errResponse
//	@Router			/currency [post]
func (h Handler) CreateCurrency(c fiber.Ctx) error {
	req, err := ParseAndValidateRequest[currencyCreateRequest](c.Request().Body())
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(errResponse{Error: errInvalidInput.Error()})
	}

	upperName := strings.ToUpper(req.Name)

	id, err := h.Currency.Create(c.Context(), domain.Currency{
		Name: upperName,
		Type: req.Type,
	})
	if err != nil {
		h.Logger.Error().Err(err).Msgf("create currency")

		var serviceErr *domain.ServiceError
		if errors.As(err, &serviceErr) {
			if serviceErr.Type == domain.Client {
				return c.Status(http.StatusBadRequest).JSON(errResponse{Error: serviceErr.Error()})
			}
		}

		return c.Status(http.StatusInternalServerError).JSON(errResponse{Error: errSomethingWentWrong.Error()})
	}

	return c.Status(http.StatusOK).JSON(createCurrencyResponse{ID: id})
}

// GetRate 		 godoc
//
//	@Summary		get currencies rate
//	@Description	get currencies rate by params
//	@Tags			currency
//	@Produce		json
//	@Param			from	query		string	true	"currency from"
//	@Param			to		query		string	true	"currency to"
//	@Param			value	query		float64	true	"currency from value"
//	@Success		200		{object}	getRateResponse
//	@Failure		400		{object}	errResponse
//	@Failure		500		{object}	errResponse
//	@Router			/currency/rate [get]
func (h Handler) GetRate(c fiber.Ctx) error {
	fromValue := c.Query(fromQueryParam)
	toValue := c.Query(toQueryParam)
	value := c.Query(valueQueryParam)

	if fromValue == "" || toValue == "" {
		return c.Status(http.StatusBadRequest).JSON(errResponse{Error: errInvalidInput.Error()})
	}

	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(errResponse{Error: errInvalidInput.Error()})
	}

	if floatValue <= 0 {
		return c.Status(http.StatusBadRequest).JSON(errResponse{Error: errInvalidInput.Error()})
	}

	fromValue = strings.ToUpper(fromValue)
	toValue = strings.ToUpper(toValue)

	rate, err := h.Currency.GetRate(c.Context(), domain.Rate{
		From:  fromValue,
		To:    toValue,
		Value: decimal.NewFromFloat(floatValue),
	})
	if err != nil {
		h.Logger.Error().Err(err).Msgf("get currency rate")

		var serviceErr *domain.ServiceError
		if errors.As(err, &serviceErr) {
			if serviceErr.Type == domain.Client {
				return c.Status(http.StatusBadRequest).JSON(errResponse{Error: serviceErr.Error()})
			}
		}

		return c.Status(http.StatusInternalServerError).JSON(errResponse{Error: errSomethingWentWrong.Error()})
	}

	return c.Status(http.StatusOK).JSON(getRateResponse{Rate: rate})
}

// ChangeCurrencyAvailability godoc
//
//	@Summary		chane currency availability
//	@Description	change currency
//	@Tags			currency
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	changeCurrencyAvailabilityRequest
//	@Failure		400	{object}	errResponse
//	@Failure		500	{object}	errResponse
//	@Router			/currency/availability [patch]
func (h Handler) ChangeCurrencyAvailability(c fiber.Ctx) error {
	req, err := ParseAndValidateRequest[changeCurrencyAvailabilityRequest](c.Request().Body())
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(errResponse{Error: errInvalidInput.Error()})
	}

	upperName := strings.ToUpper(req.Name)

	if err := h.Currency.ChangeAvailability(c.Context(), upperName, req.IsAvailable); err != nil {
		h.Logger.Error().Err(err).Msgf("change currency availability")

		var serviceErr *domain.ServiceError
		if errors.As(err, &serviceErr) {
			if serviceErr.Type == domain.Client {
				return c.Status(http.StatusBadRequest).JSON(errResponse{Error: serviceErr.Error()})
			}
		}

		return c.Status(http.StatusInternalServerError).JSON(errResponse{Error: errSomethingWentWrong.Error()})
	}

	return c.Status(http.StatusOK).JSON(defaultResponse{Success: true})
}

// GetAvailableCurrencies godoc
//
//	@Summary		get available currencies
//	@Description	get currencies
//	@Tags			currency
//	@Produce		json
//	@Success		200	{object}	getAvailableCurrenciesResponse
//	@Failure		400	{object}	errResponse
//	@Failure		500	{object}	errResponse
//	@Router			/currency/all [get]
func (h Handler) GeteCurrencies(c fiber.Ctx) error {
	currencies, err := h.Currency.GetAll(c.Context())
	if err != nil {
		h.Logger.Error().Err(err).Msgf("get available currencies")

		var serviceErr *domain.ServiceError
		if errors.As(err, &serviceErr) {
			if serviceErr.Type == domain.Client {
				return c.Status(http.StatusBadRequest).JSON(errResponse{Error: serviceErr.Error()})
			}
		}

		return c.Status(http.StatusInternalServerError).JSON(errResponse{Error: errSomethingWentWrong.Error()})
	}

	resp := make([]currency, 0, len(currencies))

	for i := range currencies {
		resp = append(resp, currencyToDto(currencies[i]))
	}

	return c.Status(http.StatusOK).JSON(getAvailableCurrenciesResponse{Currencies: resp})
}
