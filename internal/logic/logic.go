package logic

import (
	"context"

	"example.com/internal/logic/models"
	"example.com/internal/storage"
)

type Client struct {
	storage storage.Storage
}

func NewClient(storage storage.Storage) *Client {
	return &Client{
		storage: storage,
	}
}

func (c *Client) StoreTransaction(ctx context.Context, transaction models.Transaction) (models.Transaction, error) {
	return c.storage.StoreTransaction(ctx, transaction)
}

func (c *Client) GetTransaction(ctx context.Context, ID string) (models.Transaction, error) {
	return c.storage.GetTransactionByID(ctx, ID)
}
