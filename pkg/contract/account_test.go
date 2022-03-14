package contract

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAccount_Validate(t *testing.T) {
	cases := []struct {
		input    AccountRequest
		expected string
	}{
		{
			input:    AccountRequest{},
			expected: "document_number: cannot be blank.",
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.expected, func(t *testing.T) {
			assert.EqualError(t, tt.input.Validate(), tt.expected)
		})
	}
}
