
## Generate OpenAPI files


firstly install openapi generator
``go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest``

then we can use the generator using:
``~/go/bin/oapi-codegen -package wex wex/wex.yaml > ./wex/``

## Running integration tests

start the server using
``go run main.go``

then run the tests in another terminal using
``go test ./tests/integration_test.go``