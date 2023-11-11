package main

import (
	"net/http"

	"example.com/internal/server"
	"example.com/wex"
	"github.com/labstack/echo/v4"
)

func main() {
	// register our server
	handler := wex.ServerInterfaceWrapper{
		Handler: &server.Server{},
	}

	router := echo.New()
	wex.RegisterHandlers(router, &handler)

	// start the server
	err := http.ListenAndServe("localhost:6666", router.Server.Handler)
	if err != nil {
		panic(err)
	}
}
