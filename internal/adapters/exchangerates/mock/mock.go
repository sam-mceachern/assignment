package mock

import (
	"context"
	"time"
)

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

var getExchangeRateForCountry func(context.Context, string, time.Time) (float64, error)

func (c *Client) GetExchangeRateForCountry(ctx context.Context, targetCurrency string, oldestDate time.Time) (float64, error) {
	return getExchangeRateForCountry(ctx, targetCurrency, oldestDate)
}

func (C *Client) MockGetExchangeRateForCountry(mockFunc func(context.Context, string, time.Time) (float64, error)) {
	getExchangeRateForCountry = mockFunc
}
