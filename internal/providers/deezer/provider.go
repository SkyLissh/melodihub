// Package deezer provides a client for the Deezer API.
package deezer

import (
	"encoding/json"

	"github.com/skylissh/melodihub/internal/validator"
	"resty.dev/v3"
)

type Provider struct {
	client *resty.Client

	Search *SearchService
	Artist *ArtistService
}

func New(validator *validator.Validator) *Provider {
	provider := &Provider{
		client: resty.New(),
	}

	provider.client.SetBaseURL("https://api.deezer.com")

	provider.client.AddRetryConditions(func(r *resty.Response, err error) bool {
		return r.StatusCode() == 429
	})

	provider.client.AddResponseMiddleware(func(c *resty.Client, r *resty.Response) error {
		var res *Response[any]

		if err := json.Unmarshal(r.Bytes(), &res); err == nil {
			return nil
		}

		if res.Error.Code != 0 {
			return mapDeezerError(res.Error)
		}

		return nil

	})

	provider.Search = &SearchService{provider: provider, validator: validator}
	provider.Artist = &ArtistService{provider: provider, validator: validator}

	return provider
}
