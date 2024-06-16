package forex

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/alemax1/currencies-api/config"
	"github.com/alemax1/currencies-api/internal/currency/domain"
	"github.com/shopspring/decimal"
)

type CurrenciesRepo interface {
	GetCurrenciesByType(ctx context.Context, tp domain.CurrencyType) ([]domain.Currency, error)
	UpdateCurrencyValueByName(ctx context.Context, value decimal.Decimal, name string) error
}

type Forex struct {
	Client           *http.Client
	CurrenciesAPICfg config.CurrenciesAPI
}

func New(
	currenciesAPICfg config.CurrenciesAPI,
) *Forex {
	return &Forex{
		Client:           new(http.Client),
		CurrenciesAPICfg: currenciesAPICfg,
	}
}

const (
	fromQueryParam = "from"
	toQueryParam   = "to"

	usdQueryParam = "USD"

	apiKeyQueryParam = "api_key"
)

type MultiFetchResp struct {
	Results map[string]float64 `json:"results"`
}

type OneFetchResp struct {
	Result map[string]float64 `json:"result"`
}

func (f Forex) SendMultiFetchRequest(ctx context.Context, currencies []domain.Currency) ([]domain.CurrencyWithValue, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, f.CurrenciesAPICfg.FetchMultiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	query := req.URL.Query()
	query.Add(fromQueryParam, usdQueryParam)
	query.Add(apiKeyQueryParam, f.CurrenciesAPICfg.APIKey)
	var curr string
	for idx, cur := range currencies {
		if idx != len(currencies)-1 {
			curr += fmt.Sprintf("%s,", cur.Name)
			continue
		}
		curr += cur.Name
	}
	query.Add(toQueryParam, curr)
	req.URL.RawQuery = query.Encode()

	resp, err := f.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("clien do: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	var response MultiFetchResp

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("marshall body: %w", err)
	}

	currenciesResp := make([]domain.CurrencyWithValue, 0, len(currencies))

	for name, value := range response.Results {
		currenciesResp = append(currenciesResp, domain.CurrencyWithValue{
			Name:  name,
			Value: decimal.NewFromFloat(value),
		})
	}

	if len(currenciesResp) == 0 {
		return nil, domain.NewServiceError(domain.ErrNothingFound, domain.Client)
	}

	return currenciesResp, nil
}

func (f Forex) SendFetchOneRequest(ctx context.Context, currency domain.Currency) (domain.CurrencyWithValue, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, f.CurrenciesAPICfg.FetchOneURL, nil)
	if err != nil {
		return domain.CurrencyWithValue{}, fmt.Errorf("new request: %w", err)
	}
	query := req.URL.Query()
	query.Add(fromQueryParam, usdQueryParam)
	query.Add(toQueryParam, currency.Name)
	query.Add(apiKeyQueryParam, f.CurrenciesAPICfg.APIKey)
	req.URL.RawQuery = query.Encode()

	resp, err := f.Client.Do(req)
	if err != nil {
		return domain.CurrencyWithValue{}, fmt.Errorf("client do: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return domain.CurrencyWithValue{}, fmt.Errorf("read body: %w", err)
	}

	var response OneFetchResp

	if err := json.Unmarshal(body, &response); err != nil {
		return domain.CurrencyWithValue{}, fmt.Errorf("unmarshal body: %w", err)
	}

	if _, ok := response.Result[currency.Name]; !ok {
		return domain.CurrencyWithValue{}, domain.NewServiceError(domain.ErrNothingFound, domain.Client)
	}

	return domain.CurrencyWithValue{
		Name:  currency.Name,
		Value: decimal.NewFromFloat(response.Result[currency.Name]),
	}, nil
}
