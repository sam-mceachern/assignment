package fiscaldata

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"example.com/internal/logic/models"
	"example.com/util"
)

const (
	hostAddress   = "https://api.fiscaldata.treasury.gov/services/api/fiscal_service"
	exchangeRates = "v1/accounting/od/rates_of_exchange"
)

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

type Data struct {
	Data []ExchangeRate `json:"data"`
}

type ExchangeRate struct {
	ExchangeRate string `json:"exchange_rate"`
}

func (C *Client) GetExchangeRateForCountry(ctx context.Context, targetCountry string, oldestDate time.Time) (float64, error) {
	endpoint := fmt.Sprintf("%s/%s", hostAddress, exchangeRates)
	// TODO check this converts the date correctly
	endpoint = fmt.Sprintf("%s?filter=record_date:gt:%s", endpoint, oldestDate.Format(time.DateOnly))
	endpoint = fmt.Sprintf("%s,country:eq:%s", endpoint, targetCountry)

	responseBody, err := util.MakeAPICall[Data](http.MethodGet, endpoint)
	if err != nil {
		return 0.0, fmt.Errorf("failed to make api call: %w", err)
	}

	if len(responseBody.Data) == 0 {
		return 0.0, models.ErrNoExchangeRateFound
	}

	exchangeRate, err := strconv.ParseFloat(responseBody.Data[len(responseBody.Data)-1].ExchangeRate, 64)
	if err != nil {
		return 0.0, fmt.Errorf("failed unmarhsal json: %w", err)
	}

	// the dates are ordered from oldest to most recent, so we need the last result
	return exchangeRate, nil
}
