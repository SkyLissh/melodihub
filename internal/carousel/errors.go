package carousel

import (
	"errors"

	"github.com/skylissh/melodihub/internal/domain"
	"github.com/skylissh/melodihub/internal/providers/deezer"
	"github.com/skylissh/melodihub/internal/providers/lastfm"
)

func CarouselErrorFromLastfmError(err error) error {
	switch {
	case errors.Is(err, lastfm.ErrInvalidParameter),
		errors.Is(err, lastfm.ErrInvalidResource):
		return domain.InvalidParam(err)
	case errors.Is(err, lastfm.ErrRateLimited),
		errors.Is(err, lastfm.ErrTooManyRequests):
		return domain.RateLimited(domain.ProviderLastFM, err)
	case errors.Is(err, lastfm.ErrInvalidService),
		errors.Is(err, lastfm.ErrInvalidMethod),
		errors.Is(err, lastfm.ErrAuthFailed),
		errors.Is(err, lastfm.ErrInvalidFormat),
		errors.Is(err, lastfm.ErrBadAuthToken),
		errors.Is(err, lastfm.ErrOperationFailed),
		errors.Is(err, lastfm.ErrInvalidSessionKey),
		errors.Is(err, lastfm.ErrInvalidSignature),
		errors.Is(err, lastfm.ErrTemporaryServerIssues),
		errors.Is(err, lastfm.ErrInvalidUsername),
		errors.Is(err, lastfm.ErrInvalidTimestamp),
		errors.Is(err, lastfm.ErrDeletedAPIKey),
		errors.Is(err, lastfm.ErrNotEnoughContent),
		errors.Is(err, lastfm.ErrServiceUnavailable),
		errors.Is(err, lastfm.ErrUserRequired),
		errors.Is(err, lastfm.ErrSubsonicServerRequired),
		errors.Is(err, lastfm.ErrLegacyMethodDisabled),
		errors.Is(err, lastfm.ErrBadAccountScrobble),
		errors.Is(err, lastfm.ErrNonExistentError),
		errors.Is(err, lastfm.ErrRegistrationDisabled),
		errors.Is(err, lastfm.ErrAPIKeyPermissionDenied),
		errors.Is(err, lastfm.ErrSuspendedAPIKey),
		errors.Is(err, lastfm.ErrOffline),
		errors.Is(err, lastfm.ErrInvalidResponse),
		errors.Is(err, lastfm.ErrServer),
		errors.Is(err, lastfm.ErrUnknownCode):
		return domain.ProviderUnavailable(domain.ProviderLastFM, err)
	default:
		return domain.Internal(err)
	}
}

func CarouselErrorFromDeezerError(err error) error {
	switch {
	case shouldSkipCarouselItemFromDeezerError(err):
		return nil
	case errors.Is(err, deezer.ErrQuotaExceeded),
		errors.Is(err, deezer.ErrItemsLimit):
		return domain.RateLimited(domain.ProviderDeezer, err)
	case errors.Is(err, deezer.ErrPermission),
		errors.Is(err, deezer.ErrTokenInvalid),
		errors.Is(err, deezer.ErrServiceBusy),
		errors.Is(err, deezer.ErrAccountNotAllowed),
		errors.Is(err, deezer.ErrInvalidResponse),
		errors.Is(err, deezer.ErrServer),
		errors.Is(err, deezer.ErrUnknownCode):
		return domain.ProviderUnavailable(domain.ProviderDeezer, err)
	default:
		return domain.Internal(err)
	}
}

func shouldSkipCarouselItemFromDeezerError(err error) bool {
	return errors.Is(err, deezer.ErrParameterMissing) ||
		errors.Is(err, deezer.ErrInvalidParameter) ||
		errors.Is(err, deezer.ErrQueryInvalid) ||
		errors.Is(err, deezer.ErrNotFound)
}
