package worker

import (
	"context"
	"sync"
	"time"

	"github.com/alemax1/currencies-api/config"
	"github.com/alemax1/currencies-api/pkg/logger"
)

type CurrencyService interface {
	UpdateCryptoCurrencies(ctx context.Context) error
	UpdateFiatCurrencies(ctx context.Context) error
}

type Worker struct {
	CurrencyService CurrencyService
	Logger          logger.Logger
	Cfg             config.CurrenciesWorker
	stopCh          chan struct{}
	wg              sync.WaitGroup
}

func New(
	currencySvc CurrencyService,
	workerCfg config.CurrenciesWorker,
	logger logger.Logger,
) *Worker {
	return &Worker{
		CurrencyService: currencySvc,
		Logger:          logger,
		Cfg:             workerCfg,
		stopCh:          make(chan struct{}),
		wg:              sync.WaitGroup{},
	}
}

func (w *Worker) Run() {
	w.wg.Add(1)
	ticker := time.NewTicker(w.Cfg.IterationTimeout)
	for {
		select {
		case <-w.stopCh:
			w.Logger.Info().Msg("worker successfully shuted down")
			w.wg.Done()
			return
		default:
			<-ticker.C
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

			if err := w.CurrencyService.UpdateFiatCurrencies(ctx); err != nil {
				w.Logger.Error().Err(err).Msgf("update fiat currencies")
			}
			w.Logger.Info().Msg("fiat currencies successfully updated")

			if err := w.CurrencyService.UpdateCryptoCurrencies(ctx); err != nil {
				w.Logger.Error().Err(err).Msgf("update fiat currencies")
			}
			w.Logger.Info().Msg("crypto currencies successfully updated")

			cancel()
		}
	}
}

func (w *Worker) Stop() {
	w.stopCh <- struct{}{}
	w.wg.Wait()
}
