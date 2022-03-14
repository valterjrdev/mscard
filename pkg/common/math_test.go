package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIntToNegative(t *testing.T) {

	cases := []struct {
		input    int64
		expected int64
	}{
		{
			input:    1,
			expected: 1,
		},
		{
			input:    -1,
			expected: 1,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run("", func(t *testing.T) {
			assert.Equal(t, tt.expected, Abs(tt.input))
		})
	}

}
