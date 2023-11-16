package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"example.com/internal/logic/models"
	"example.com/wex"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/rs/zerolog/log"
)

func (s *Server) PostStoreTransaction(ctx echo.Context) error {
	// extract the request
	req, err := getRequestStruct[wex.StoreTransactionRequest](ctx, s.apiRouter)
	if err != nil {
		log.Err(err).Msgf("failed to get request struct")
		return writeErrorResponse(ctx.Response(), http.StatusBadRequest, err.Error())
	}

	err = validatePurchaseAmount(req.PurchaseAmountUSD)
	if err != nil {
		log.Err(err).Msgf("failed validate purchase amount")
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
	return models.Transaction{
		Description:     req.Description,
		TransactionDate: req.TransactionDate.Time,
		PurchaseAmount:  float64(req.PurchaseAmountUSD),
	}, nil
}

func TransactionToStoreTransactionResponse(transaction models.Transaction) wex.StoreTransactionResponse {
	return wex.StoreTransactionResponse{
		Id:          *transaction.ID,
		Description: transaction.Description,
		TransactionDate: openapi_types.Date{
			Time: transaction.TransactionDate,
		},
		PurchaseAmountUSD: float32(transaction.PurchaseAmount),
	}
}

// sadly the openapi validator is unable to check if a number is rounded to 2 decimal places.
// this is a know issue in the library: https://github.com/getkin/kin-openapi/issues/817
// so here we are validating this field manually
func validatePurchaseAmount(purchaseAmount float32) error {
	amountStr := strconv.FormatFloat(float64(purchaseAmount), 'f', -1, 32)
	amountSplit := strings.Split(amountStr, ".")
	if len(amountSplit) < 2 {
		return nil
	}

	if len(amountSplit[1]) > 2 {
		return fmt.Errorf("purchase amountUSD has too many decimal places: %d", len(amountSplit[1]))
	}

	return nil
}
