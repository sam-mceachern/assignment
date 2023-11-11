package server

import (
	"encoding/json"
	"io"

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
	}

	var req wex.StoreTransactionRequest
	err = json.Unmarshal(data, &req)
	if err != nil {
		log.Info().Msgf("failed to unmarshal transaction: %v", err)
	}

	return nil
}

// (POST /storeTransaction)
func (s *Server) PostStoreTransaction(ctx echo.Context) error {
	log.Print("beer")
	return nil
}
