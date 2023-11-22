package storage

import (
	"context"

	"example.com/internal/logic/models"
)

// Storage defines a generic interface for storing transaction information
type Storage interface {
	StoreTransaction(context.Context, models.Transaction) (models.Transaction, error)
	GetTransactionByID(context.Context, string) (models.Transaction, error)
}
