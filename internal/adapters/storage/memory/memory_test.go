package memory

import (
	"context"
	"errors"
	"testing"
	"time"

	"example.com/internal/logic/models"
	"github.com/stretchr/testify/assert"
)

func TestStoreTransaction(t *testing.T) {
	tests := []struct {
		name             string
		inputTransaction models.Transaction
		expectedOutput   models.Transaction
	}{
		{
			name: "add transaction",
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
		},
	}

	client := NewClient()
	for _, test := range tests {
		output, err := client.StoreTransaction(context.Background(), test.inputTransaction)
		assert.NoError(t, err, test.name)

		assert.Equal(t, test.expectedOutput.Description, output.Description, test.name)
		assert.Equal(t, test.expectedOutput.PurchaseAmount, output.PurchaseAmount, test.name)
		assert.Equal(t, test.expectedOutput.TransactionDate, output.TransactionDate, test.name)
		assert.NotEqual(t, "", output.ID)
	}
}

func TestGetTransaction(t *testing.T) {
	client := NewClient()
	transaction, err := client.StoreTransaction(context.Background(), models.Transaction{
		Description:     "hello",
		TransactionDate: time.Unix(123, 0),
		PurchaseAmount:  123,
	})
	assert.NoError(t, err)

	tests := []struct {
		name           string
		ID             string
		expectedOutput models.Transaction
		expectedErr    error
	}{
		{
			name: "get transaction",
			ID:   *transaction.ID,
			expectedOutput: models.Transaction{
				Description:     "hello",
				TransactionDate: time.Unix(123, 0),
				PurchaseAmount:  123,
			},
		},
		{
			name:           "no transaction found",
			ID:             "",
			expectedOutput: models.Transaction{},
			expectedErr:    errors.New("could not find result"),
		},
	}

	for _, test := range tests {
		output, err := client.GetTransactionByID(context.Background(), test.ID)
		if test.expectedErr != nil {
			assert.Equal(t, test.expectedErr, err)
			continue
		}

		assert.NoError(t, err, test.name)
		assert.Equal(t, test.expectedOutput.Description, output.Description, test.name)
		assert.Equal(t, test.expectedOutput.PurchaseAmount, output.PurchaseAmount, test.name)
		assert.Equal(t, test.expectedOutput.TransactionDate, output.TransactionDate, test.name)
		assert.NotEqual(t, test.ID, output.ID)
	}
}
