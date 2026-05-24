package deezer

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResponseUnmarshalDecodesSuccessfulPayloadIntoGenericData(t *testing.T) {
	var response Response[Artist]

	err := json.Unmarshal([]byte(`{
		"id": 27,
		"name": "Daft Punk",
		"link": "https://www.deezer.com/artist/27",
		"picture_medium": "https://e-cdns-images.dzcdn.net/images/artist/test/250x250-000000-80-0-0.jpg"
	}`), &response)

	require.NoError(t, err)
	require.Nil(t, response.Error)
	require.NotNil(t, response.Data)
	assert.Equal(t, 27, response.Data.ID)
	assert.Equal(t, "Daft Punk", response.Data.Name)
}

func TestResponseUnmarshalDecodesNestedDeezerErrorPayload(t *testing.T) {
	var response Response[Artist]

	err := json.Unmarshal([]byte(`{
		"error": {
			"type": "DataException",
			"message": "no data",
			"code": 800
		}
	}`), &response)

	require.NoError(t, err)
	require.NotNil(t, response.Error)
	require.Nil(t, response.Data)
	assert.Equal(t, DataNotFoundCode, response.Error.Code)
	assert.Equal(t, "DataException", response.Error.Kind)
	assert.Equal(t, "no data", response.Error.Message)
}
