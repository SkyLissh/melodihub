package carousel

import (
	"errors"

	"github.com/skylissh/melodihub/internal/domain"
	"github.com/skylissh/melodihub/internal/providers/deezer"
	"github.com/skylissh/melodihub/internal/providers/lastfm"
)

func CarouselErrorFromLastfmError(err error) error {
	var lastfmErr lastfm.LastfmError
	if errors.As(err, &lastfmErr) {
		switch {
		case errors.Is(lastfmErr.Kind, lastfm.ErrInvalidParameter),
			errors.Is(lastfmErr.Kind, lastfm.ErrInvalidResource):
			return domain.InvalidParam(err)
		case errors.Is(lastfmErr.Kind, lastfm.ErrRateLimited),
			errors.Is(lastfmErr.Kind, lastfm.ErrTooManyRequests):
			return domain.RateLimited(domain.ProviderLastFM, err)
		default:
			return domain.ProviderUnavailable(domain.ProviderLastFM, err)
		}
	}
	return domain.Internal(err)
}

func CarouselErrorFromDeezerError(err error) error {
	var deezerErr deezer.DeezerError
	if errors.As(err, &deezerErr) {
		switch {
		case shouldSkipCarouselItemFromDeezerKind(deezerErr.Kind):
			return nil
		case errors.Is(deezerErr.Kind, deezer.ErrQuotaExceeded),
			errors.Is(deezerErr.Kind, deezer.ErrItemsLimit):
			return domain.RateLimited(domain.ProviderDeezer, err)
		default:
			return domain.ProviderUnavailable(domain.ProviderDeezer, err)
		}
	}
	return domain.Internal(err)
}

func shouldSkipCarouselItemFromDeezerKind(kind error) bool {
	return errors.Is(kind, deezer.ErrParameterMissing) ||
		errors.Is(kind, deezer.ErrInvalidParameter) ||
		errors.Is(kind, deezer.ErrQueryInvalid) ||
		errors.Is(kind, deezer.ErrNotFound)
}
