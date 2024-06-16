package config

type CurrenciesAPI struct {
	APIKey        string
	FetchMultiURL string
	FetchOneURL   string
}

func newCurrenciesAPI() CurrenciesAPI {
	return CurrenciesAPI{
		APIKey:        getDefaultEnv("CURRENCIES_API_KEY", ""),
		FetchMultiURL: getDefaultEnv("CURRENCIES_API_FETCH_MULTI_URL", ""),
		FetchOneURL:   getDefaultEnv("CURRENCIES_API_FETCH_ONE_URL", ""),
	}
}
