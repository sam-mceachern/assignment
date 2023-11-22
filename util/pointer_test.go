package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToPtr(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected interface{}
	}{
		{
			name:     "simple",
			value:    123,
			expected: 123,
		},
		{
			name:     "simple",
			value:    "123",
			expected: "123",
		},
	}

	for _, test := range tests {
		actualValue := ToPtr(test.value)
		assert.Equal(t, test.expected, *actualValue)
	}
}
