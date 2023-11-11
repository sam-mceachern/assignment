package models

import "time"

type Transaction struct {
	ID              *string
	Description     string
	TransactionDate time.Time
	PurchaseAmount  float64
}
