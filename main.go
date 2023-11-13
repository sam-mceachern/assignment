package main

import (
	"context"
	"net/http"
	"os"

	"example.com/internal/logic"
	"example.com/internal/server"
	"example.com/internal/storage/memory"
	"example.com/wex"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func main() {

	loader := &openapi3.Loader{Context: context.Background(), IsExternalRefsAllowed: true}
	doc, err := loader.LoadFromFile("./wex/wex.yaml")
	if err != nil {
		log.Err(err).Msg("failed to load openapi file")
		os.Exit(1)
	}

	err = doc.Validate(context.Background())
	if err != nil {
		log.Err(err).Msg("failed validate document")
		os.Exit(1)
	}

	apirouter, err := gorillamux.NewRouter(doc)
	if err != nil {
		log.Err(err).Msg("failed to create router")
		os.Exit(1)
	}

	// comppose the three layers of our API

	// 1. firstly our storage layer
	storage := memory.NewClient()

	// 2. then our logic layer
	logicClient := logic.NewClient(storage)

	// 3. finally our transport layer
	server := server.NewServer(logicClient, apirouter)

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
}
