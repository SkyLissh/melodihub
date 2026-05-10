package domain

import "errors"

var (
	ErrInvalidProvider = errors.New("invalid provider")
)

type Provider string

const (
	ProviderDeezer Provider = "deezer"
	ProviderLastFM Provider = "lastfm"
)

func ParseProvider(s string) (Provider, error) {
	switch s {
	case "deezer":
		return ProviderDeezer, nil
	case "lastfm":
		return ProviderLastFM, nil
	default:
		return "", ErrInvalidProvider
	}
}

func (p Provider) String() string {
	return string(p)
}

func (p Provider) IsValid() bool {
	return p == ProviderDeezer || p == ProviderLastFM
}
