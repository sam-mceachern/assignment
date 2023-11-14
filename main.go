package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"example.com/internal/logic"
	"example.com/internal/server"
	"example.com/internal/storage/memory"
	"example.com/wex"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func main() {
	// load our validation router, this helps us to validate incoming API requests
	apiRouter, err := createValidationRouter()
	if err != nil {
		log.Err(err).Msg("failed to create validation router")
		os.Exit(1)
	}

	// compose the three layers of our API
	// 1. firstly our storage layer
	storage := memory.NewClient()

	// 2. then our logic layer
	logicClient := logic.NewClient(storage)

	// 3. finally our transport layer
	server := server.NewServer(logicClient, apiRouter)

	// register our server
	handler := wex.ServerInterfaceWrapper{
		Handler: server,
	}
	router := echo.New()
	wex.RegisterHandlers(router, &handler)

	// start the server
	err = http.ListenAndServe("localhost:5555", router.Server.Handler)
	if err != nil {
		panic(err)
	}

	log.Info().Msg("api stopped")
}

func createValidationRouter() (routers.Router, error) {
	loader := &openapi3.Loader{Context: context.Background(), IsExternalRefsAllowed: true}
	doc, err := loader.LoadFromFile("./wex/wex.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to load api file: %v", err)
	}

	err = doc.Validate(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to validate api spec: %v", err)
	}

	apirouter, err := gorillamux.NewRouter(doc)
	if err != nil {
		return nil, fmt.Errorf("failed to create router: %v", err)
	}

	return apirouter, nil
}
