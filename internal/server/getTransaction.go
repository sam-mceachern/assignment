package server

import (
	"net/http"

	"example.com/wex"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func (s *Server) PostGetTransaction(ctx echo.Context) error {
	// extract the request
	req, err := getRequestStruct[wex.GetTransactionRequest](ctx, s.apiRouter)
	if err != nil {
		log.Err(err).Msgf("failed to get request struct")
		return writeErrorResponse(ctx.Response(), http.StatusInternalServerError, "internal error")
	}

	// get the transaction
	transaction, err := s.logicClient.GetTransaction(ctx.Request().Context(), req.ID)
	if err != nil {
		log.Err(err).Msgf("failed to store transaction")
		return writeErrorResponse(ctx.Response(), http.StatusInternalServerError, "internal error")
	}

	// write the transaction to the response
	return writeResponse(ctx.Response(), http.StatusCreated, TransactionToStoreTransactionResponse(transaction))
}