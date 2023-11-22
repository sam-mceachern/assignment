package mock

import (
	"context"

	"example.com/internal/logic/models"
)

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

var storeTransaction func(context.Context, models.Transaction) (models.Transaction, error)
var getTransaction func(context.Context, string) (models.Transaction, error)

func (C *Client) StoreTransaction(ctx context.Context, transaction models.Transaction) (models.Transaction, error) {
	return storeTransaction(ctx, transaction)
}

func (C *Client) GetTransactionByID(ctx context.Context, id string) (models.Transaction, error) {
	return getTransaction(ctx, id)
}

func (C *Client) MockStoreTransaction(mockFunc func(context.Context, models.Transaction) (models.Transaction, error)) {
	storeTransaction = mockFunc
}

func (C *Client) MockGetTransaction(mockFunc func(context.Context, string) (models.Transaction, error)) {
	getTransaction = mockFunc
}
