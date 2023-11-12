package memory

import (
	"context"
	"encoding/json"
	"fmt"

	"example.com/internal/logic/models"
	"github.com/google/uuid"
)

type Client struct {
	storage map[string][]byte
}

func NewClient() *Client {
	return &Client{
		storage: make(map[string][]byte),
	}
}

func (c *Client) StoreTransaction(ctx context.Context, transaction models.Transaction) (models.Transaction, error) {
	id := generateID()
	transaction.ID = &id

	data, err := json.Marshal(transaction)
	if err != nil {
		return models.Transaction{}, fmt.Errorf("failed to marshal transaction to json: %v", err)
	}

	c.storage[id] = data

	return transaction, nil
}

func (c *Client) GetTransactionByID(ctx context.Context, id string) (models.Transaction, error) {
	return models.Transaction{}, nil
}

func generateID() string {
	return uuid.New().String()
}
