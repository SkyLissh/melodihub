// Package core provides the core configuration for the application.
package core

import (
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type Env struct {
	LastfmAPIKey string `validate:"required"`
	ValkeyAddr   string `validate:"required"`
}

func NewEnv() (*Env, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	env := &Env{
		LastfmAPIKey: os.Getenv("LASTFM_API_KEY"),
		ValkeyAddr:   os.Getenv("VALKEY_ADDR"),
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.Struct(env)
	if err != nil {
		return nil, err
	}

	return env, nil
}
