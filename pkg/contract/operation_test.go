package contract

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContractOperationType_Validate_Error(t *testing.T) {
	cases := []struct {
		input    OperationRequest
		expected string
	}{
		{
			input:    OperationRequest{},
			expected: "debit: cannot be blank; description: cannot be blank.",
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.expected, func(t *testing.T) {
			assert.EqualError(t, tt.input.Validate(), tt.expected)
		})
	}
}
