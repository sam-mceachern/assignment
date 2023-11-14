package server

import (
	"fmt"
	"net/http"
	"time"

	"example.com/internal/logic/models"
	"example.com/wex"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func (s *Server) PostStoreTransaction(ctx echo.Context) error {
	// extract the request
	req, err := getRequestStruct[wex.StoreTransactionRequest](ctx, s.apiRouter)
	if err != nil {
		log.Err(err).Msgf("failed to get request struct")
		return writeErrorResponse(ctx.Response(), http.StatusBadRequest, err.Error())
	}

	// convert the request to our logic layer type
	transaction, err := StoreTransactionRequestToTransaction(req)
	if err != nil {
		log.Err(err).Msgf("failed to convert storeTransactionRequest into transaction")
		return writeErrorResponse(ctx.Response(), http.StatusBadRequest, err.Error())
	}

	// store the transaction
	transaction, err = s.logicClient.StoreTransaction(ctx.Request().Context(), transaction)
	if err != nil {
		log.Err(err).Msgf("failed to store transaction")
		return writeErrorResponse(ctx.Response(), http.StatusInternalServerError, "internal error")
	}

	// write the result in the response
	return writeResponse(ctx.Response(), http.StatusCreated, TransactionToStoreTransactionResponse(transaction))
}

func StoreTransactionRequestToTransaction(req wex.StoreTransactionRequest) (models.Transaction, error) {
	transactionDate, err := time.Parse(time.DateOnly, req.TransactionDate)
	if err != nil {
		return models.Transaction{}, fmt.Errorf("failed to parse transaction date: %v", err)
	}

	return models.Transaction{
		Description:     req.Description,
		TransactionDate: transactionDate,
		PurchaseAmount:  float64(req.PurchaseAmount),
	}, nil
}

func TransactionToStoreTransactionResponse(transaction models.Transaction) wex.StoreTransactionResponse {
	return wex.StoreTransactionResponse{
		ID:              *transaction.ID,
		Description:     transaction.Description,
		TransactionDate: transaction.TransactionDate.String(),
		PurchaseAmount:  float32(transaction.PurchaseAmount),
	}
}
