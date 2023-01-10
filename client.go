package cbr

import (
	"net/http"
	"time"

	"github.com/shopspring/decimal"
)

// fetchFunction is a function that mimics http.Get() method
type fetchFunction func(url string) (resp *http.Response, err error)

// Client is a currency rates service client... what else?
type Client interface {
	GetRate(string, time.Time) (float64, error)
	GetRateDecimal(string, time.Time) (decimal.Decimal, error)
	GetRateString(string, time.Time) (string, error)
	GetCurrencyInfo(string, time.Time) (Currency, error)
	SetFetchFunction(fetchFunction)
}

type client struct {
	fetch fetchFunction
}

// Returns currency rate in float64
func (s client) GetRate(currency string, t time.Time) (float64, error) {
	return getRate(currency, t, s.fetch)
}

// Returns currency rate in Decimal
//
// Rationale: https://pkg.go.dev/github.com/shopspring/decimal - FAQ section
func (s client) GetRateDecimal(currency string, t time.Time) (decimal.Decimal, error) {
	return getRateDecimal(currency, t, s.fetch)
}

// Returns currency rate string with dot as decimal separator
func (s client) GetRateString(currency string, t time.Time) (string, error) {
	return getRateString(currency, t, s.fetch)
}

// Returns currency struct
func (s client) GetCurrencyInfo(currency string, t time.Time) (Currency, error) {
	return getCurrency(currency, t, s.fetch)
}

func (s client) SetFetchFunction(f fetchFunction) {
	s.fetch = f
}

// NewClient creates a new rates service instance
func NewClient() Client {
	return client{http.Get}
}
