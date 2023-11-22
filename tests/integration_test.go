package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"example.com/util"
	"example.com/wex"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/assert"
)

type requestType int

const (
	storeTransaction requestType = iota
	getTransaction
)

const (
	apiEndpoint = "http://localhost:5555"
)

type testCase struct {
	name        string
	requestType requestType
	requestBody interface{}

	expectedResponse interface{}
	exptectedStatus  int
	expectError      bool
	expectedError    wex.ErrorResponse
}

func (t *testCase) setGetRequestID(id string) {
	if req, ok := t.requestBody.(wex.GetTransactionRequest); ok {
		req.Id = id
		t.requestBody = req
	}
}

var testCases = []testCase{
	{
		name:        "simple - store transaction",
		requestType: storeTransaction,
		requestBody: wex.StoreTransactionRequest{
			Description:       "hello",
			PurchaseAmountUSD: 123.35,
			TransactionDate: openapi_types.Date{
				Time: time.Date(2020, 2, 2, 0, 0, 0, 0, time.UTC),
			},
		},
		expectedResponse: wex.StoreTransactionResponse{
			Description:       "hello",
			PurchaseAmountUSD: 123.35,
			TransactionDate: openapi_types.Date{
				Time: time.Date(2020, 2, 2, 0, 0, 0, 0, time.UTC),
			},
		},
		exptectedStatus: http.StatusOK,
	},
	{
		name:        "simple - get transaction",
		requestType: getTransaction,
		requestBody: wex.GetTransactionRequest{
			Currency: "Australia",
		},
		expectedResponse: wex.GetTransactionResponse{
			Description:       "hello",
			PurchaseAmountUSD: 123.35,
			TransactionDate: openapi_types.Date{
				Time: time.Date(2020, 2, 2, 0, 0, 0, 0, time.UTC),
			},
			PurchaseAmountTargetCurrency: 123,
		},
		exptectedStatus: http.StatusOK,
	},
	{
		name:        "too many decimal places - store transaction",
		requestType: storeTransaction,
		requestBody: wex.StoreTransactionRequest{
			Description:       "hello",
			PurchaseAmountUSD: 123.0982345,
			TransactionDate: openapi_types.Date{
				Time: time.Date(2020, 2, 2, 0, 0, 0, 0, time.UTC),
			},
		},
		expectedResponse: wex.StoreTransactionResponse{},
		expectError:      true,
		expectedError: wex.ErrorResponse{
			Description: "number has too many decimal places: 6",
		},
		exptectedStatus: http.StatusBadRequest,
	},
	{
		name:        "no description - store transaction",
		requestType: storeTransaction,
		requestBody: wex.StoreTransactionRequest{
			Description:       "",
			PurchaseAmountUSD: 123.09,
			TransactionDate: openapi_types.Date{
				Time: time.Date(2020, 2, 2, 0, 0, 0, 0, time.UTC),
			},
		},
		expectedResponse: wex.StoreTransactionResponse{},
		expectError:      true,
		expectedError: wex.ErrorResponse{
			Description: "failed to validate request: request body has an error: doesn't match schema: Error at \"/description\": minimum string length is 1\nSchema:\n  {\n    \"maxLength\": 50,\n    \"minLength\": 1,\n    \"type\": \"string\"\n  }\n\nValue:\n  \"\"\n",
		},
		exptectedStatus: http.StatusBadRequest,
	},
	{
		name:        "no currency exists - get transaction",
		requestType: getTransaction,
		requestBody: wex.GetTransactionRequest{
			Currency: "Antarctica",
		},
		expectedResponse: wex.GetTransactionResponse{},
		expectError:      true,
		expectedError: wex.ErrorResponse{
			Description: "could find exchange rate for: 'Antarctica'",
		},
		exptectedStatus: http.StatusInternalServerError,
	},
}

func TestThing(t *testing.T) {
	var lastTransactionID string

	for _, testCase := range testCases {
		switch testCase.requestType {
		case storeTransaction:
			resp, errResp, err := makeRequest[wex.StoreTransactionResponse](testCase, "/storeTransaction")
			if err != nil {
				assert.Fail(t, err.Error())
			}

			if errResp == nil {
				lastTransactionID = resp.Id
			}

			assertStoreTransactionResponse(t, testCase, resp, errResp)
		case getTransaction:
			testCase.setGetRequestID(lastTransactionID)
			resp, errResp, err := makeRequest[wex.GetTransactionResponse](testCase, "/getTransaction")
			if err != nil {
				assert.Fail(t, err.Error())
			}

			assertGetTransactionResponse(t, testCase, lastTransactionID, resp, errResp)
		}
	}
}

func assertStoreTransactionResponse(t *testing.T, test testCase, response wex.StoreTransactionResponse, errResp *wex.ErrorResponse) {
	if test.expectError {
		assertError(t, test, errResp)
	}

	expectedResponse, ok := test.expectedResponse.(wex.StoreTransactionResponse)
	if !ok {
		assert.Fail(t, "failed to convert response to StoreTransactionResponse", test.name)
	}

	assert.Equal(t, expectedResponse.Description, response.Description, test.name)
	assert.Equal(t, expectedResponse.PurchaseAmountUSD, response.PurchaseAmountUSD, test.name)
	assert.Equal(t, expectedResponse.TransactionDate, response.TransactionDate, test.name)
}

func assertGetTransactionResponse(t *testing.T, test testCase, lastTransactionID string, response wex.GetTransactionResponse, errResp *wex.ErrorResponse) {
	if test.expectError {
		assertError(t, test, errResp)
		return
	}

	expectedResponse, ok := test.expectedResponse.(wex.GetTransactionResponse)
	if !ok {
		assert.Fail(t, "failed to convert response to GetTransactionResponse", test.name)
	}

	assert.Equal(t, expectedResponse.Description, response.Description, test.name)
	assert.Equal(t, expectedResponse.PurchaseAmountUSD, response.PurchaseAmountUSD, test.name)
	assert.Equal(t, expectedResponse.TransactionDate, response.TransactionDate, test.name)
	assert.Equal(t, lastTransactionID, response.Id, test.name)

	// these two cant be reliably compared, as exchange rates change over time
	// assert.Equal(t, response.ExchangeRate, expectedResponse.ExchangeRate)
	// assert.Equal(t, response.PurchaseAmountTargetCurrency, expectedResponse.PurchaseAmountTargetCurrency)

	assert.NoError(t, util.CheckNumberIsRoundedTo(response.PurchaseAmountTargetCurrency, 2), test.name)
}

func assertError(t *testing.T, test testCase, errResp *wex.ErrorResponse) {
	assert.Equal(t, test.expectedError.Description, errResp.Description, test.name)
}

func makeRequest[O any](test testCase, path string) (O, *wex.ErrorResponse, error) {
	var output O
	endpoint := apiEndpoint + path

	requestData, err := json.Marshal(test.requestBody)
	if err != nil {
		return output, nil, fmt.Errorf("failed to marshal request")
	}

	reader := bytes.NewReader(requestData)
	responseData, status, err := util.MakeAPICall(http.MethodPost, endpoint, reader)
	if err != nil {
		return output, nil, err
	}

	if status > 299 {
		var errResp wex.ErrorResponse
		err = json.Unmarshal(responseData, &errResp)
		if err != nil {
			return output, nil, fmt.Errorf("failed to unmarshal error response: '%w'", err)
		}

		return output, &errResp, nil
	}

	err = json.Unmarshal(responseData, &output)
	if err != nil {
		return output, nil, fmt.Errorf("failed to unmarshal response: '%w'", err)
	}

	return output, nil, nil
}
