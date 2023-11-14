package server

import (
	"fmt"
	"net/http"

	"example.com/internal/logic/models"
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
		if err == models.ErrCouldNotFindResult {
			return writeErrorResponse(ctx.Response(), http.StatusNotFound, fmt.Sprintf("could not find result for: '%s'", req.ID))
		}
		return writeErrorResponse(ctx.Response(), http.StatusInternalServerError, "internal error")
	}

	// write the transaction to the response
	return writeResponse(ctx.Response(), http.StatusCreated, TransactionToStoreTransactionResponse(transaction))
}
