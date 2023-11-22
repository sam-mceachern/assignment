package util

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckNumberIsRoundedTo(t *testing.T) {
	tests := []struct {
		name          string
		inputNumber   float32
		decimalPlaces int
		expectedError error
	}{
		{
			name:          "two decimal places",
			inputNumber:   123.321,
			decimalPlaces: 2,
			expectedError: errors.New("number has too many decimal places: 3"),
		},
		{
			name:          "number passes",
			inputNumber:   123.3214,
			decimalPlaces: 4,
			expectedError: nil,
		},
		{
			name:          "number passes - no decimal places",
			inputNumber:   123,
			decimalPlaces: 4,
			expectedError: nil,
		},
	}

	for _, test := range tests {
		err := CheckNumberIsRoundedTo(test.inputNumber, test.decimalPlaces)
		assert.Equal(t, err, test.expectedError)
	}
}
