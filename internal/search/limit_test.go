package search

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseLimitAcceptsEmptyDefaultAndValidBounds(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Limit
	}{
		{name: "empty uses default", input: "", want: DefaultLimit},
		{name: "minimum", input: "1", want: Limit(1)},
		{name: "normal value", input: "25", want: Limit(25)},
		{name: "maximum", input: "100", want: Limit(100)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limit, err := ParseLimit(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.want, limit)
		})
	}
}

func TestParseLimitRejectsInvalidValues(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  error
	}{
		{name: "zero", input: "0", want: ErrLimitTooSmall},
		{name: "negative", input: "-1", want: ErrLimitInvalid},
		{name: "not a number", input: "many", want: ErrLimitInvalid},
		{name: "decimal", input: "1.5", want: ErrLimitInvalid},
		{name: "above maximum", input: "101", want: ErrLimitTooLarge},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limit, err := ParseLimit(tt.input)

			require.Error(t, err)
			assert.True(t, errors.Is(err, tt.want))
			assert.Zero(t, limit)
		})
	}
}
