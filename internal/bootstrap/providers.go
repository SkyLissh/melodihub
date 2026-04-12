package bootstrap

import (
	"github.com/skylissh/melodihub/internal/core"
	"github.com/skylissh/melodihub/internal/providers/deezer"
	"github.com/skylissh/melodihub/internal/providers/lastfm"
)

type Providers struct {
	Deezer *deezer.Provider
	Lastfm *lastfm.LastfmProvider
}

func NewProviders(env *core.Env) *Providers {
	deezerProvider := deezer.New()
	lastfmProvider := lastfm.New(env.LastfmAPIKey)

	return &Providers{
		Deezer: deezerProvider,
		Lastfm: lastfmProvider,
	}
}
