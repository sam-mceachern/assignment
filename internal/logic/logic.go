package logic

import (
	"context"
	"fmt"
	"math"

	"example.com/internal/adapters/exchangerates"
	"example.com/internal/adapters/storage"
	"example.com/internal/logic/models"
)

type Client struct {
	storage       storage.Storage
	exchangeRates exchangerates.ExchangeRates
}

func NewClient(storage storage.Storage, exhangeRate exchangerates.ExchangeRates) *Client {
	return &Client{
		storage:       storage,
		exchangeRates: exhangeRate,
	}
}

func (c *Client) StoreTransaction(ctx context.Context, transaction models.Transaction) (models.Transaction, error) {
	return c.storage.StoreTransaction(ctx, transaction)
}

func (c *Client) GetTransaction(ctx context.Context, ID, country string) (models.Transaction, float64, float64, error) {
	transaction, err := c.storage.GetTransactionByID(ctx, ID)
	if err != nil {
		return models.Transaction{}, 0, 0, fmt.Errorf("failed to get transaction: %w", err)
	}

	sixMonthsFromPurchaseDate := transaction.TransactionDate.AddDate(0, -6, 0)
	exchangeRate, err := c.exchangeRates.GetExchangeRateForCountry(ctx, country, sixMonthsFromPurchaseDate)
	if err != nil {
		return models.Transaction{}, 0, 0, fmt.Errorf("failed to get exchange rate: %w", err)
	}

	amountTagetCurrency := math.Round(transaction.PurchaseAmount*exchangeRate*100) / 100

	return transaction, exchangeRate, amountTagetCurrency, nil
}
