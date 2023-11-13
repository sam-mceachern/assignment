package server

import (
	"example.com/internal/logic"
	"github.com/getkin/kin-openapi/routers"
)

type Server struct {
	logicClient *logic.Client
	apiRouter   routers.Router
}

func NewServer(logicClient *logic.Client, apiRouter routers.Router) *Server {
	return &Server{
		logicClient: logicClient,
		apiRouter:   apiRouter,
	}
}
