package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseIDWithValidDeezerID(t *testing.T) {
	id, err := ParseID("deezer:10583405")

	require.NoError(t, err)
	assert.Equal(t, ProviderDeezer, id.Provider())
	assert.Equal(t, 10583405, id.Id())
	assert.Equal(t, "deezer:10583405", id.String())
}

func TestParseIDWithValidLastfmID(t *testing.T) {
	id, err := ParseID("lastfm:42")

	require.NoError(t, err)
	assert.Equal(t, ProviderLastFM, id.Provider())
	assert.Equal(t, 42, id.Id())
	assert.Equal(t, "lastfm:42", id.String())
}

func TestParseIDRejectsEmptyString(t *testing.T) {
	_, err := ParseID("")

	require.ErrorIs(t, err, ErrEmptyId)
}

func TestParseIDRejectsMissingColon(t *testing.T) {
	_, err := ParseID("deezer123")

	require.ErrorIs(t, err, ErrInvalidId)
}

func TestParseIDRejectsNonNumericID(t *testing.T) {
	_, err := ParseID("deezer:abc")

	require.ErrorIs(t, err, ErrInvalidId)
}

func TestParseIDRejectsUnknownProvider(t *testing.T) {
	_, err := ParseID("spotify:123")

	require.ErrorIs(t, err, ErrInvalidProvider)
}

func TestParseIDRejectsNegativeID(t *testing.T) {
	_, err := ParseID("deezer:-1")

	require.ErrorIs(t, err, ErrInvalidId)
}

func TestParseIDRejectsZeroID(t *testing.T) {
	id, err := ParseID("deezer:0")

	require.NoError(t, err)
	assert.Equal(t, 0, id.Id())
}

func TestNewDeezerID(t *testing.T) {
	id, err := NewDeezerID(10583405)

	require.NoError(t, err)
	assert.Equal(t, ProviderDeezer, id.Provider())
	assert.Equal(t, 10583405, id.Id())
	assert.Equal(t, "deezer:10583405", id.String())
}

func TestNewLastFMID(t *testing.T) {
	id, err := NewLastFMID(7)

	require.NoError(t, err)
	assert.Equal(t, ProviderLastFM, id.Provider())
	assert.Equal(t, 7, id.Id())
	assert.Equal(t, "lastfm:7", id.String())
}

func TestNewIDRejectsNegativeID(t *testing.T) {
	_, err := NewID("deezer", -1)

	require.ErrorIs(t, err, ErrInvalidId)
}

func TestNewIDRejectsUnknownProvider(t *testing.T) {
	_, err := NewID("spotify", 1)

	require.ErrorIs(t, err, ErrInvalidProvider)
}

func TestParseProviderValid(t *testing.T) {
	tests := []struct {
		input    string
		expected Provider
	}{
		{"deezer", ProviderDeezer},
		{"lastfm", ProviderLastFM},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			p, err := ParseProvider(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, p)
		})
	}
}

func TestParseProviderRejectsUnknown(t *testing.T) {
	_, err := ParseProvider("spotify")
	require.ErrorIs(t, err, ErrInvalidProvider)
}

func TestProviderIsValid(t *testing.T) {
	assert.True(t, ProviderDeezer.IsValid())
	assert.True(t, ProviderLastFM.IsValid())
	assert.False(t, Provider("spotify").IsValid())
}

func TestProviderString(t *testing.T) {
	assert.Equal(t, "deezer", ProviderDeezer.String())
	assert.Equal(t, "lastfm", ProviderLastFM.String())
}