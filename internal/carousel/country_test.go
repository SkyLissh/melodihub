package carousel

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseCountryDefaultsWhenEmpty(t *testing.T) {
	country, err := ParseCountry("")

	require.NoError(t, err)
	assert.Equal(t, DefaultCountry, country)
	assert.Equal(t, "united states", country.String())
}

func TestParseCountryNormalizesWhitespaceAndCasing(t *testing.T) {
	country, err := ParseCountry("  United   States  ")

	require.NoError(t, err)
	assert.Equal(t, Country("united states"), country)
	assert.Equal(t, "united states", country.String())
}

func TestParseCountryRejectsInvalidCharacters(t *testing.T) {
	country, err := ParseCountry("mexico!")

	require.ErrorIs(t, err, ErrCountryInvalid)
	assert.Empty(t, country)
}
