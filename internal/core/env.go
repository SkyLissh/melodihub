// Package core provides the core configuration for the application.
package core

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Env struct {
	LastfmAPIKey string
	ValkeyAddr   string
}

func NewEnv() (*Env, error) {
	if err := godotenv.Load(); err != nil {
		log.Warn().Err(err).Msg("failed to load .env file")
	}

	env := &Env{
		LastfmAPIKey: os.Getenv("LASTFM_API_KEY"),
		ValkeyAddr:   os.Getenv("VALKEY_ADDR"),
	}

	if env.LastfmAPIKey == "" || env.ValkeyAddr == "" {
		return nil, os.ErrInvalid
	}

	return env, nil
}
