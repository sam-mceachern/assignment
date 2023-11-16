package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"example.com/wex"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/labstack/echo/v4"
)

func getRequestStruct[T any](ctx echo.Context, router routers.Router) (T, error) {
	var reqStruct T

	request := ctx.Request()
	defer request.Body.Close()

	err := validateRequest(ctx.Request().Context(), request, router)
	if err != nil {
		return reqStruct, err
	}

	data, err := io.ReadAll(request.Body)
	if err != nil {
		return reqStruct, fmt.Errorf("failed to read request body: %w", err)
	}

	err = json.Unmarshal(data, &reqStruct)
	if err != nil {
		return reqStruct, fmt.Errorf("failed to unmarshal transaction: %w", err)
	}

	return reqStruct, nil
}

func validateRequest(ctx context.Context, req *http.Request, router routers.Router) error {
	route, pathParams, err := router.FindRoute(req)
	if err != nil {
		return fmt.Errorf("failed to find router: %w", err)
	}

	// Validate request
	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request:    req,
		PathParams: pathParams,
		Route:      route,
	}

	err = openapi3filter.ValidateRequest(ctx, requestValidationInput)
	if err != nil {
		return fmt.Errorf("failed to validate request: %w", err)
	}

	return nil
}

func writeResponse(w http.ResponseWriter, status int, response interface{}) error {
	w.WriteHeader(status)

	w.Header().Set("Content-Type", "application/json")
	responseData, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	_, err = w.Write(responseData)
	if err != nil {
		return fmt.Errorf("failed to write data: %w", err)
	}

	return nil
}

func writeErrorResponse(w http.ResponseWriter, status int, message string) error {
	w.WriteHeader(status)

	w.Header().Set("Content-Type", "application/json")
	errorResponse := wex.ErrorResponse{
		Description: message,
	}
	errorData, err := json.Marshal(errorResponse)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	_, err = w.Write(errorData)
	if err != nil {
		return fmt.Errorf("failed to write data: %w", err)
	}

	return nil
}
