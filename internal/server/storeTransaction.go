package server

import (
	"net/http"

	"example.com/internal/logic/models"
	"example.com/util"
	"example.com/wex"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/rs/zerolog/log"
)

func (s *Server) PostStoreTransaction(ctx echo.Context) error {
	log.Info().Msg("PostStoreTransaction")

	// extract the request
	req, err := getRequestStruct[wex.StoreTransactionRequest](ctx, s.apiRouter)
	if err != nil {
		log.Err(err).Msgf("failed to get request struct")
		return writeErrorResponse(ctx.Response(), http.StatusBadRequest, err.Error())
	}

	err = util.CheckNumberIsRoundedTo(req.PurchaseAmountUSD, 2)
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
