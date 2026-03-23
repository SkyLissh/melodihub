// Package deezer provides a client for the Deezer API.
package deezer

import (
	"resty.dev/v3"
)

type Provider struct {
	client *resty.Client

	Search SearchService
	Artist ArtistService
}

func New() *Provider {
	provider := &Provider{
		client: resty.New(),
	}

	provider.client.SetBaseURL("https://api.deezer.com")

	provider.client.AddRetryConditions(func(r *resty.Response, err error) bool {
		return r.StatusCode() == 429
	})

	provider.Search = &searchService{provider: provider}
	provider.Artist = &artistService{provider: provider}

	return provider
}
