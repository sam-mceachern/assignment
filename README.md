
## Running API
run api using
``go run main.go``

run unit tests using:
``go test ./...``

## example curls to use

### store transaction
``curl --location --request POST 'localhost:5555/storeTransaction' \
--header 'Content-Type: application/json' \
--data-raw '{
    "description": "i am a transaction",
    "transactionDate": "2023-01-02",
    "purchaseAmountUSD": 563.03
}'``

## get transaction
``curl --location --request POST 'localhost:5555/getTransaction' \
--header 'Content-Type: application/json' \
--data-raw '{
    "Id": "<target_id_here>",
    "currency": "Australia"
}'``

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

