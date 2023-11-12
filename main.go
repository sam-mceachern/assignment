package main

import (
	"net/http"

	"example.com/internal/logic"
	"example.com/internal/server"
	"example.com/internal/storage/memory"
	"example.com/wex"
	"github.com/labstack/echo/v4"
)

func main() {
	// comppose the three layers of our API

	// 1. firstly our storage layer
	storage := memory.NewClient()

	// 2. then our logic layer
	logicClient := logic.NewClient(storage)

	// 3. finally our transport layer
	server := server.NewServer(logicClient)

	// register our server
	handler := wex.ServerInterfaceWrapper{
		Handler: server,
	}

	router := echo.New()
	wex.RegisterHandlers(router, &handler)

	// start the server
	err := http.ListenAndServe("localhost:6666", router.Server.Handler)
	if err != nil {
		panic(err)
	}

}
