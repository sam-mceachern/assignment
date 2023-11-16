package server

import (
	"errors"
	"fmt"
	"net/http"

	"example.com/internal/logic/models"
	"example.com/wex"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/rs/zerolog/log"
)

func (s *Server) PostGetTransaction(ctx echo.Context) error {
	// extract the request
	req, err := getRequestStruct[wex.GetTransactionRequest](ctx, s.apiRouter)
	if err != nil {
		log.Err(err).Msgf("failed to get request struct")
		return writeErrorResponse(ctx.Response(), http.StatusBadRequest, err.Error())
	}

	// get the transaction
	transaction, exchangeRate, amount, err := s.logicClient.GetTransaction(ctx.Request().Context(), req.Id, req.Currency)
	if err != nil {
		log.Err(err).Msgf("failed to get transaction")
		if errors.Is(err, models.ErrCouldNotFindResult) {
			return writeErrorResponse(ctx.Response(), http.StatusNotFound, fmt.Sprintf("could not find result for: '%s'", req.Id))
		}
		if errors.Is(err, models.ErrNoExchangeRateFound) {
			return writeErrorResponse(ctx.Response(), http.StatusNotFound, fmt.Sprintf("could find exchange rate for: '%s'", req.Currency))
		}
		return writeErrorResponse(ctx.Response(), http.StatusInternalServerError, "internal error")
	}

	// write the transaction to the response
	return writeResponse(ctx.Response(), http.StatusCreated, TransactionToGetTransactionResponse(transaction, exchangeRate, amount))
}

func TransactionToGetTransactionResponse(transaction models.Transaction, exchangeRate, amount float64) wex.GetTransactionResponse {
	return wex.GetTransactionResponse{
		Id:          *transaction.ID,
		Description: transaction.Description,
		TransactionDate: openapi_types.Date{
			Time: transaction.TransactionDate,
		},
		PurchaseAmountUSD:            float32(transaction.PurchaseAmount),
		ExchangeRate:                 float32(exchangeRate),
		PurchaseAmountTargetCurrency: float32(amount),
	}
}
