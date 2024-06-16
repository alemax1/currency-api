package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/alemax1/currencies-api/internal/currency/domain"
	"github.com/shopspring/decimal"
)

type Currency struct {
	*DBExecutor
}

func NewCurrency(executor *DBExecutor) *Currency {
	return &Currency{
		DBExecutor: executor,
	}
}

func (c Currency) AddEmptyCurrency(ctx context.Context, currency domain.Currency) (int64, error) {
	var id int64

	currency.ValueUSD = decimal.Zero
	currency.IsAvailable = false

	if err := c.db.QueryRowContext(ctx,
		"INSERT INTO currencies(type, name, is_available, value_usd) VALUES($1, $2, $3, $4) RETURNING id",
		currency.Type,
		currency.Name,
		currency.IsAvailable,
		currency.ValueUSD,
	).Scan(
		&id,
	); err != nil {
		if strings.Contains(err.Error(), duplicateErrorCode) {
			return 0, domain.NewServiceError(domain.ErrDuplicateValue, domain.Client)
		}

		return 0, err
	}

	return id, nil
}

func (c Currency) UpdateCurrencyAvailability(ctx context.Context, name string, isAvailable bool) error {
	result, err := c.db.ExecContext(ctx,
		"UPDATE currencies SET is_available=$1, updated_at=CURRENT_TIMESTAMP WHERE name=$2",
		isAvailable,
		name,
	)
	if err != nil {
		return newExecContextErr(err)
	}

	updatedRows, err := result.RowsAffected()
	if err != nil {
		return newUpdatedRowsErr(err)
	}
	if updatedRows == 0 {
		return domain.NewServiceError(domain.ErrNothingFound, domain.Client)
	}

	return nil
}

func (c Currency) UpdateCurrencyByName(ctx context.Context, currency domain.CurrencyUpdateData) error {
	result, err := c.db.ExecContext(ctx,
		"UPDATE currencies SET value_usd=$1, is_available=$2, updated_at=CURRENT_TIMESTAMP WHERE name=$3",
		currency.ValueUSD,
		currency.IsAvailable,
		currency.Name,
	)
	if err != nil {
		return newExecContextErr(err)
	}

	updatedRows, err := result.RowsAffected()
	if err != nil {
		return newUpdatedRowsErr(err)
	}
	if updatedRows == 0 {
		return domain.NewServiceError(domain.ErrNothingUpdated, domain.Client)
	}

	return nil
}

func (c Currency) GetCurrenciesByType(ctx context.Context, tp domain.CurrencyType) ([]domain.Currency, error) {
	rows, err := c.db.QueryContext(ctx,
		"SELECT id, name, type, value_usd, is_available FROM currencies WHERE type=$1",
		tp,
	)
	if err != nil {
		return nil, newQueryErr(err)
	}
	defer rows.Close()

	var currencies []domain.Currency

	for rows.Next() {
		var currency domain.Currency

		if err := rows.Scan(
			&currency.ID,
			&currency.Name,
			&currency.Type,
			&currency.ValueUSD,
			&currency.IsAvailable,
		); err != nil {
			return nil, newScanErr(err)
		}

		currencies = append(currencies, currency)
	}
	if rows.Err() != nil {
		return nil, newRowsErr(err)
	}

	return currencies, nil
}

func (c Currency) GetCurrency(ctx context.Context, name string) (domain.Currency, error) {
	var currency domain.Currency

	if err := c.db.QueryRowContext(ctx,
		"SELECT id, name, type, value_usd, is_available FROM currencies WHERE name=$1", name,
	).Scan(
		&currency.ID,
		&currency.Name,
		&currency.Type,
		&currency.ValueUSD,
		&currency.IsAvailable,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Currency{}, domain.NewServiceError(domain.ErrNothingFound, domain.Client)
		}
	}

	return currency, nil
}

func (c Currency) GetAll(ctx context.Context) ([]domain.Currency, error) {
	rows, err := c.db.QueryContext(ctx,
		"SELECT id, name, type, value_usd, is_available FROM currencies",
	)
	if err != nil {
		return nil, newQueryErr(err)
	}
	defer rows.Close()

	var currencies []domain.Currency

	for rows.Next() {
		var currency domain.Currency

		if err := rows.Scan(
			&currency.ID,
			&currency.Name,
			&currency.Type,
			&currency.ValueUSD,
			&currency.IsAvailable,
		); err != nil {
			return nil, newScanErr(err)
		}

		currencies = append(currencies, currency)
	}
	if rows.Err() != nil {
		return nil, newRowsErr(err)
	}

	return currencies, nil
}
