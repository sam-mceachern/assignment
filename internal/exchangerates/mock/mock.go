package mock

import (
	"context"
	"time"
)

type Client struct {
	out float64
	err error
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) GetExchangeRateForCountry(ctx context.Context, targetCurrency string, oldestDate time.Time) (float64, error) {
	return c.out, c.err
}

func (c *Client) GetExchangeRateForCountryWillReturn(out float64, err error) {
	c.out = out
	c.err = err
}
