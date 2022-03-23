package contract

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContractTransaction_Validate_Error(t *testing.T) {
	cases := []struct {
		description string
		input       TransactionRequest
		expected    string
	}{
		{
			description: "fields required",
			input:       TransactionRequest{},
			expected:    "account_id: cannot be blank; amount: cannot be blank; operation_id: cannot be blank.",
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			assert.EqualError(t, tt.input.Validate(), tt.expected)
		})
	}
}
