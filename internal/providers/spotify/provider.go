// Package spotify provides a Spotify API provider
package spotify

import (
	"net/http"
	"slices"
	"strconv"
	"sync"
	"time"

	"resty.dev/v3"
)

var retryStatusCodes = []int{
	http.StatusTooManyRequests,
	http.StatusInternalServerError,
	http.StatusUnauthorized,
}

type SpotifyProvider struct {
	ClientID     string
	ClientSecret string

	client *resty.Client

	Artist SpotifyArtist
	Album  SpotifyAlbum
	Search SpotifySearch

	token *Token
	mu    sync.Mutex
}

func (s *SpotifyProvider) ensureToken() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var token *Token

	if s.token != nil && !s.token.IsExpired() {
		return nil
	}

	client := resty.New()
	defer client.Close()

	_, err := client.R().
		SetFormData(map[string]string{
			"grant_type":    "client-credentials",
			"client_id":     s.ClientID,
			"client_secret": s.ClientSecret,
		}).
		SetResult(&token).
		Post("https://accounts.spotify.com/api/token")

	if err != nil {
		return err
	}

	s.token = token

	return nil
}

func (s *SpotifyProvider) invalidateToken() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.token = nil
}

func New(clientID string, clientSecret string) *SpotifyProvider {
	provider := &SpotifyProvider{
		ClientID:     clientID,
		ClientSecret: clientSecret,

		client: resty.New(),
	}

	provider.client.SetBaseURL("https://api.spotify.com/v1/")

	provider.client.AddRetryConditions(
		func(res *resty.Response, err error) bool {
			if err != nil {
				return false
			}

			status := res.StatusCode()

			if status == http.StatusUnauthorized {
				provider.invalidateToken()
			}

			return slices.Contains(retryStatusCodes, status)
		},
	)

	provider.client.SetRetryStrategy(
		func(res *resty.Response, err error) (time.Duration, error) {
			if err != nil {
				return 0, nil
			}

			if res.StatusCode() == http.StatusTooManyRequests {
				retryAfter := res.Header().Get("Retry-After")

				if sec, errConv := strconv.Atoi(retryAfter); errConv == nil {
					return time.Duration(sec) * time.Second, nil
				}
			}

			return 0, nil
		},
	)

	provider.client.AddRequestMiddleware(
		func(ctx *resty.Client, req *resty.Request) error {
			if err := provider.ensureToken(); err != nil {
				return err
			}

			req.SetAuthToken(provider.token.AccessToken)

			return nil
		},
	)

	provider.Artist = &spotifyArtist{provider: provider}
	provider.Album = &spotifyAlbum{provider: provider}
	provider.Search = &spotifySearch{provider: provider}

	return provider
}
