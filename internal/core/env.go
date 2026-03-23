// Package core provides the core configuration for the application.
package core

import (
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type Env struct {
	LastfmAPIKey        string `validate:"required"`
	SpotifyClientID     string `validate:"required"`
	SpotifyClientSecret string `validate:"required"`
	ValkeyAddr          string `validate:"required"`
}

func NewEnv() *Env {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
		os.Exit(1)
	}

	env := &Env{
		LastfmAPIKey:        os.Getenv("LASTFM_API_KEY"),
		SpotifyClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		SpotifyClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		ValkeyAddr:          os.Getenv("VALKEY_ADDR"),
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.Struct(env)
	if err != nil {
		log.Fatal("Error validating .env file", err.Error())
		os.Exit(1)
	}

	return env
}
