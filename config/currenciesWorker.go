package config

import "time"

type CurrenciesWorker struct {
	IterationTimeout time.Duration
}

func newCurrenciesWorker() CurrenciesWorker {
	return CurrenciesWorker{
		IterationTimeout: getDefaultDurationEnv("CURRENCIES_WORKER_ITERATION_TIMEOUT", 1*time.Minute),
	}
}
