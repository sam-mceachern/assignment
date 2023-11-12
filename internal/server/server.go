package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"example.com/internal/logic"
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
	// extract the request
	defer ctx.Request().Body.Close()
	data, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		log.Info().Msgf("failed to read request body data: %v", err)
		return writeResponse(ctx.Response(), http.StatusBadRequest, "failed to read request body")
	}

	var req wex.StoreTransactionRequest
	err = json.Unmarshal(data, &req)
	if err != nil {
		log.Info().Msgf("failed to unmarshal transaction: %v", err)
		return writeResponse(ctx.Response(), http.StatusBadRequest, "failed to unmarshal data into request body")
	}

	return nil
}

// (POST /storeTransaction)
func (s *Server) PostStoreTransaction(ctx echo.Context) error {
	log.Print("beer")
	return nil
}

func writeResponse(w http.ResponseWriter, status int, message string) error {
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
