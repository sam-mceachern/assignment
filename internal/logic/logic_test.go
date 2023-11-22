package logic

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	exchangeRateMock "example.com/internal/adapters/exchangerates/mock"
	storageMock "example.com/internal/adapters/storage/mock"
	"example.com/internal/logic/models"
	"example.com/util"
	"github.com/stretchr/testify/assert"
)

func TestStoreTransaction(t *testing.T) {
	storageMock := storageMock.NewClient()
	exchangeRateMock := exchangeRateMock.NewClient()
	logicClient := NewClient(storageMock, exchangeRateMock)
	tests := []struct {
		name             string
		inputTransaction models.Transaction
		expectedOutput   models.Transaction
		expectedError    error
		mockAction       func()
	}{
		{
			name: "store transaction",
			inputTransaction: models.Transaction{
				Description:     "hello",
				TransactionDate: time.Unix(123, 0),
				PurchaseAmount:  123,
			},
			expectedOutput: models.Transaction{
				Description:     "hello",
				TransactionDate: time.Unix(123, 0),
				PurchaseAmount:  123,
			},
			mockAction: func() {
				storageMock.MockStoreTransaction(func(ctx context.Context, transaction models.Transaction) (models.Transaction, error) {
					return models.Transaction{
						Description:     "hello",
						TransactionDate: time.Unix(123, 0),
						PurchaseAmount:  123,
					}, nil
				})
			},
		},
		{
			name: "store transaction - downstream returns error",
			inputTransaction: models.Transaction{
				Description:     "hello",
				TransactionDate: time.Unix(123, 0),
				PurchaseAmount:  123,
			},
			expectedOutput: models.Transaction{},
			expectedError:  errors.New("hello"),
			mockAction: func() {
				storageMock.MockStoreTransaction(func(ctx context.Context, transaction models.Transaction) (models.Transaction, error) {
					return models.Transaction{}, errors.New("hello")
				})
			},
		},
	}

	for _, test := range tests {
		test.mockAction()

		transaction, err := logicClient.StoreTransaction(context.Background(), test.inputTransaction)
		if test.expectedError != nil {
			assert.Equal(t, test.expectedError, err, test.name)
		}

		assert.Equal(t, test.expectedOutput, transaction, test.name)
	}
}

func TestGetTransaction(t *testing.T) {
	storageMock := storageMock.NewClient()
	exchangeRateMock := exchangeRateMock.NewClient()
	logicClient := NewClient(storageMock, exchangeRateMock)
	tests := []struct {
		name    string
		ID      string
		country string

		exchangeRate         float64
		amountTargetCurrency float64
		expectedOutput       models.Transaction
		expectedError        error
		mockAction           func()
	}{
		{
			name: "get transaction",
			ID:   "321",
			expectedOutput: models.Transaction{
				ID:              util.ToPtr("321"),
				Description:     "hello",
				TransactionDate: time.Unix(123, 0),
				PurchaseAmount:  123,
			},
			country:              "cool town",
			exchangeRate:         1.5,
			amountTargetCurrency: 184.5,
			mockAction: func() {
				storageMock.MockGetTransaction(func(ctx context.Context, id string) (models.Transaction, error) {
					return models.Transaction{
						ID:              &id,
						Description:     "hello",
						TransactionDate: time.Unix(123, 0),
						PurchaseAmount:  123,
					}, nil
				})

				exchangeRateMock.MockGetExchangeRateForCountry(func(ctx context.Context, targetCountry string, oldestDate time.Time) (float64, error) {
					return 1.5, nil
				})
			},
		},
		{
			name:           "get transaction - storage returns error",
			ID:             "321",
			expectedOutput: models.Transaction{},
			country:        "cool town",
			mockAction: func() {
				storageMock.MockGetTransaction(func(ctx context.Context, id string) (models.Transaction, error) {
					return models.Transaction{}, errors.New("oh no")
				})
			},
			expectedError: fmt.Errorf("failed to get transaction: %w", errors.New("oh no")),
		},
		{
			name:           "get transaction - exchange rate api returns errorr",
			ID:             "321",
			expectedOutput: models.Transaction{},
			country:        "cool town",
			mockAction: func() {
				storageMock.MockGetTransaction(func(ctx context.Context, id string) (models.Transaction, error) {
					return models.Transaction{
						ID:              &id,
						Description:     "hello",
						TransactionDate: time.Unix(123, 0),
						PurchaseAmount:  123,
					}, nil
				})

				exchangeRateMock.MockGetExchangeRateForCountry(func(ctx context.Context, targetCountry string, oldestDate time.Time) (float64, error) {
					return 0.0, errors.New("oh jeez")
				})
			},
			expectedError: fmt.Errorf("failed to get exchange rate: %w", errors.New("oh jeez")),
		},
	}

	for _, test := range tests {
		test.mockAction()

		transaction, exchangeRate, targetCurrencyAmount, err := logicClient.GetTransaction(context.Background(), test.ID, test.country)
		if test.expectedError != nil {
			assert.Equal(t, test.expectedError, err, test.name)
		}

		assert.Equal(t, test.expectedOutput, transaction, test.name)
		assert.Equal(t, test.exchangeRate, exchangeRate, test.name)
		assert.Equal(t, test.amountTargetCurrency, targetCurrencyAmount, test.name)
	}
}
