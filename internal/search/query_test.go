package search

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseQueryAcceptsLettersDigitsSpacesAndUnicodeLetters(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "ascii letters", input: "muse"},
		{name: "letters digits and spaces", input: "muse 1997"},
		{name: "japanese kanji", input: "坂本龍一"},
		{name: "japanese kana", input: "宇多田ヒカル"},
		{name: "korean hangul", input: "아이유"},
		{name: "chinese han", input: "周杰倫"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, err := ParseQuery(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.input, query.String())
		})
	}
}

func TestParseQueryRejectsEmptyOrBlankInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "empty", input: ""},
		{name: "spaces only", input: "   "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, err := ParseQuery(tt.input)

			require.Error(t, err)
			assert.True(t, errors.Is(err, ErrInvalidQuery))
			assert.Empty(t, query)
		})
	}
}

func TestParseQueryRejectsSymbolCharacters(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "opening brace", input: "muse{"},
		{name: "closing bracket", input: "muse]"},
		{name: "backslash", input: `muse\`},
		{name: "plus sign", input: "muse+"},
		{name: "equals sign", input: "muse="},
		{name: "minus sign", input: "muse-"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, err := ParseQuery(tt.input)

			require.Error(t, err)
			assert.True(t, errors.Is(err, ErrInvalidChar))
			assert.Empty(t, query)
		})
	}
}
