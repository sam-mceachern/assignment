package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"example.com/internal/logic"
	"example.com/internal/logic/models"
	"example.com/wex"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type Server struct {
	logicClient *logic.Client
}

func NewServer(logicClient *logic.Client) *Server {
	return &Server{logicClient: logicClient}
}

// (POST /getTransaction)
func (s *Server) PostGetTransaction(ctx echo.Context) error {
	return nil
}

func StoreTransactionRequestToTransaction(req wex.StoreTransactionRequest) (models.Transaction, error) {
	transactionDate, err := time.Parse(time.RFC3339, req.TransactionDate)
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
		Id:              *transaction.ID,
		Description:     transaction.Description,
		TransactionDate: transaction.TransactionDate.String(),
		PurchaseAmount:  float32(transaction.PurchaseAmount),
	}
}

// (POST /storeTransaction)
func (s *Server) PostStoreTransaction(ctx echo.Context) error {
	// extract the request
	defer ctx.Request().Body.Close()
	data, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		log.Info().Msgf("failed to read request body data: %v", err)
		return writeErrorResponse(ctx.Response(), http.StatusBadRequest, "failed to read request body")
	}

	var req wex.StoreTransactionRequest
	err = json.Unmarshal(data, &req)
	if err != nil {
		log.Info().Msgf("failed to unmarshal transaction: %v", err)
		return writeErrorResponse(ctx.Response(), http.StatusBadRequest, "failed to unmarshal data into request body")
	}

	transaction, err := StoreTransactionRequestToTransaction(req)
	if err != nil {
		log.Info().Msgf("failed to convert storeTransactionRequest into transaction: %v", err)
		return writeErrorResponse(ctx.Response(), http.StatusBadRequest, err.Error())
	}

	transaction, err = s.logicClient.StoreTransaction(ctx.Request().Context(), transaction)
	if err != nil {
		log.Info().Msgf("failed to storee transaction: %v", err)
		return writeErrorResponse(ctx.Response(), http.StatusInternalServerError, "internal error")
	}

	return writeResponse(ctx.Response(), http.StatusCreated, TransactionToStoreTransactionResponse(transaction))
}

func writeResponse(w http.ResponseWriter, status int, response interface{}) error {
	w.WriteHeader(status)

	w.Header().Set("Content-Type", "application/json")
	responseData, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal: %v", err)
	}

	_, err = w.Write(responseData)
	if err != nil {
		return fmt.Errorf("failed to write data: %v", err)
	}

	return nil
}

func writeErrorResponse(w http.ResponseWriter, status int, message string) error {
	w.WriteHeader(status)

	w.Header().Set("Content-Type", "application/json")
	errorResponse := wex.ErrorResponse{
		Message: message,
	}
	errorData, err := json.Marshal(errorResponse)
	if err != nil {
		return fmt.Errorf("failed to marshal: %v", err)
	}

	_, err = w.Write(errorData)
	if err != nil {
		return fmt.Errorf("failed to write data: %v", err)
	}

	return nil
}
