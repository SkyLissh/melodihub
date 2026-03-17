// Package lastfm provides a client for the Last.fm API.
package lastfm

import (
	"encoding/json"
	"net/url"
	"strings"

	"resty.dev/v3"
)

type LastfmProvider struct {
	client *resty.Client

	APIKey string

	Geo    LastfmGeo
	Artist LastfmArtist
	Album  LastfmAlbum
	Track  LastfmTrack
}

func New(apiKey string) *LastfmProvider {
	l := &LastfmProvider{client: resty.New(), APIKey: apiKey}

	l.client.SetBaseURL("http://ws.audioscrobbler.com/2.0/")

	l.client.AddRequestMiddleware(func(ctx *resty.Client, req *resty.Request) error {
		newURL, err := url.Parse(req.URL)
		if err != nil {
			return err
		}

		endpoint := strings.ReplaceAll(newURL.Path, "/2.0/", "")

		newURL.Path = ""

		q := newURL.Query()
		q.Set("api_key", l.APIKey)
		q.Set("format", "json")
		q.Set("method", endpoint)
		newURL.RawQuery = q.Encode()

		req.SetURL(newURL.String())

		return nil
	})

	l.client.AddResponseMiddleware(func(ctx *resty.Client, res *resty.Response) error {
		apiErr := &responseError{}
		if err := json.Unmarshal(res.Bytes(), apiErr); err != nil {
			return nil // non-JSON response, not an API error
		}

		if apiErr.Error != 0 {
			return newResponseError(apiErr.Error, apiErr.Message)
		}

		return nil
	})

	l.Track = &lastfmTrack{provider: l}
	l.Album = &lastfmAlbum{provider: l}
	l.Artist = &lastfmArtist{provider: l}
	l.Geo = &lastfmGeo{provider: l}

	return l
}
