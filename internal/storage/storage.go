package storage

import (
	"context"

	"example.com/internal/logic/models"
)

type Storage interface {
	StoreTransaction(context.Context, models.Transaction) (models.Transaction, error)
	GetTransactionByID(context.Context, string) (models.Transaction, error)
}
