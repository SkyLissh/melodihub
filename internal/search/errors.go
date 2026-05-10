package search

import (
	"errors"

	"github.com/skylissh/melodihub/internal/domain"
	"github.com/skylissh/melodihub/internal/providers/deezer"
)

func SearchErrorFromDeezerError(err error) error {
	switch {
	case errors.Is(err, deezer.ErrParameterMissing),
		errors.Is(err, deezer.ErrInvalidParameter),
		errors.Is(err, deezer.ErrQueryInvalid):
		return domain.InvalidParam(err)
	case errors.Is(err, deezer.ErrNotFound):
		return domain.NotFound(err)
	case errors.Is(err, deezer.ErrQuotaExceeded),
		errors.Is(err, deezer.ErrItemsLimit):
		return domain.RateLimited(domain.ProviderDeezer, err)
	case errors.Is(err, deezer.ErrServiceBusy),
		errors.Is(err, deezer.ErrInvalidResponse),
		errors.Is(err, deezer.ErrServer),
		errors.Is(err, deezer.ErrUnknownCode):
		return domain.ProviderUnavailable(domain.ProviderDeezer, err)
	default:
		return domain.Internal(err)
	}
}
