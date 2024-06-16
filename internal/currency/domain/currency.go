package domain

import "github.com/shopspring/decimal"

type CurrencyType string

const (
	Fiat   CurrencyType = "fiat"
	Crypto CurrencyType = "crypto"
)

type Currency struct {
	ID          int64
	Name        string
	Type        CurrencyType
	ValueUSD    decimal.Decimal
	IsAvailable bool
}

type CurrencyUpdateData struct {
	Name        string
	ValueUSD    decimal.Decimal
	IsAvailable bool
}

type CurrencyWithValue struct {
	Name  string
	Value decimal.Decimal
}

type Rate struct {
	From  string
	To    string
	Value decimal.Decimal
}
