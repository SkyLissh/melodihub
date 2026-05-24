package bootstrap

import (
	"github.com/skylissh/melodihub/internal/core"
	"github.com/skylissh/melodihub/internal/providers/deezer"
	"github.com/skylissh/melodihub/internal/providers/lastfm"
	"github.com/skylissh/melodihub/internal/validator"
)

type Providers struct {
	Deezer *deezer.Provider
	Lastfm *lastfm.LastfmProvider
}

func NewProviders(env *core.Env, validator *validator.Validator) *Providers {
	deezerProvider := deezer.New(validator)
	lastfmProvider := lastfm.New(env.LastfmAPIKey)

	return &Providers{
		Deezer: deezerProvider,
		Lastfm: lastfmProvider,
	}
}
