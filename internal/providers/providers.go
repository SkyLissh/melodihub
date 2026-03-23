// Package providers define a struct that holds all the providers
package providers

import (
	"github.com/skylissh/melodihub/internal/providers/deezer"
	"github.com/skylissh/melodihub/internal/providers/lastfm"
	"github.com/skylissh/melodihub/internal/providers/spotify"
)

type Providers struct {
	Spotify *spotify.SpotifyProvider
	Lastfm  *lastfm.LastfmProvider
	Deezer  *deezer.Provider
}
