package exchangerates

import (
	"context"
	"time"
)

// Storage defines a generic interface for storing transaction information
type ExchangeRates interface {
	GetExchangeRateForCountry(ctx context.Context, targetCurrency string, oldestDate time.Time) (float64, error)
}
