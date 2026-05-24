package search

import (
	"errors"

	"github.com/skylissh/melodihub/internal/domain"
	"github.com/skylissh/melodihub/internal/providers/deezer"
)

func SearchErrorFromDeezerError(err error) error {
	var deezerErr deezer.DeezerError
	if errors.As(err, &deezerErr) {
		switch {
		case errors.Is(deezerErr.Kind, deezer.ErrParameterMissing),
			errors.Is(deezerErr.Kind, deezer.ErrInvalidParameter),
			errors.Is(deezerErr.Kind, deezer.ErrQueryInvalid):
			return domain.InvalidParam(err)
		case errors.Is(deezerErr.Kind, deezer.ErrNotFound):
			return domain.NotFound(err)
		case errors.Is(deezerErr.Kind, deezer.ErrQuotaExceeded),
			errors.Is(deezerErr.Kind, deezer.ErrItemsLimit):
			return domain.RateLimited(domain.ProviderDeezer, err)
		default:
			return domain.ProviderUnavailable(domain.ProviderDeezer, err)
		}
	}
	return domain.Internal(err)
}
