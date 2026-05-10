package lastfm

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResponseUnmarshalDecodesSuccessfulPayloadIntoGenericData(t *testing.T) {
	var response Response[struct {
		TopTracks TopTracks `json:"toptracks"`
	}]

	err := json.Unmarshal([]byte(`{
		"toptracks": {
			"track": [
				{
					"name": "Hysteria",
					"duration": "227",
					"mbid": "",
					"url": "https://www.last.fm/music/Muse/_/Hysteria",
					"artist": {
						"name": "Muse",
						"mbid": "",
						"url": "https://www.last.fm/music/Muse"
					},
					"listeners": "100"
				}
			]
		}
	}`), &response)

	require.NoError(t, err)
	require.Nil(t, response.Error)
	require.Len(t, response.Data.TopTracks.Tracks, 1)
	assert.Equal(t, "Hysteria", response.Data.TopTracks.Tracks[0].Name)
}

func TestResponseUnmarshalDecodesFlatLastfmErrorPayload(t *testing.T) {
	var response Response[struct {
		TopTracks TopTracks `json:"toptracks"`
	}]

	err := json.Unmarshal([]byte(`{"error":6,"message":"Invalid parameters"}`), &response)

	require.NoError(t, err)
	require.NotNil(t, response.Error)
	assert.Equal(t, InvalidParametersCode, response.Error.Code)
	assert.Equal(t, "Invalid parameters", response.Error.Message)
}
