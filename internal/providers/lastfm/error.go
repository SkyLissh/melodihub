package lastfm

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidService         = errors.New("invalid service")
	ErrInvalidMethod          = errors.New("invalid method")
	ErrAuthFailed             = errors.New("authentication failed")
	ErrInvalidFormat          = errors.New("invalid format")
	ErrBadAuthToken           = errors.New("bad authentication token")
	ErrInvalidParameter       = errors.New("invalid parameter")
	ErrInvalidResource        = errors.New("invalid resource")
	ErrOperationFailed        = errors.New("operation failed")
	ErrInvalidSessionKey      = errors.New("invalid session key")
	ErrInvalidSignature       = errors.New("invalid method signature")
	ErrTemporaryServerIssues  = errors.New("temporary server issues")
	ErrInvalidUsername        = errors.New("invalid username")
	ErrInvalidTimestamp       = errors.New("invalid timestamp")
	ErrDeletedAPIKey          = errors.New("deleted api key")
	ErrNotEnoughContent       = errors.New("not enough content")
	ErrServiceUnavailable     = errors.New("service unavailable")
	ErrUserRequired           = errors.New("user required")
	ErrSubsonicServerRequired = errors.New("subsonic server required")
	ErrLegacyMethodDisabled   = errors.New("legacy method disabled")
	ErrBadAccountScrobble     = errors.New("bad account scrobble")
	ErrNonExistentError       = errors.New("non-existent error code")
	ErrRegistrationDisabled   = errors.New("registration disabled")
	ErrTooManyRequests        = errors.New("too many requests")
	ErrAPIKeyPermissionDenied = errors.New("api key permission denied")
	ErrSuspendedAPIKey        = errors.New("suspended api key")
	ErrRateLimited            = errors.New("rate limited")
	ErrOffline                = errors.New("offline")
	ErrInvalidResponse        = errors.New("invalid response")
	ErrServer                 = errors.New("server error")
	ErrUnknownCode            = errors.New("unknown error code")
)

type LastfmError struct {
	Kind error
	Msg  string
}

func (e LastfmError) Error() string {
	if e.Msg == "" {
		return e.Kind.Error()
	}

	return fmt.Sprintf("%v: %s", e.Kind, e.Msg)
}

func (e LastfmError) Unwrap() error {
	return e.Kind
}

func NewLastfmError(kind error, msg string) LastfmError {
	return LastfmError{
		Kind: kind,
		Msg:  msg,
	}
}

func ServerError(err error) error {
	return NewLastfmError(ErrServer, err.Error())
}

func InvalidResponse(err error) error {
	return NewLastfmError(ErrInvalidResponse, err.Error())
}

func UnknownCode(code ErrorCode, msg string) error {
	return NewLastfmError(ErrUnknownCode, fmt.Sprintf("code=%d msg=%s", code, msg))
}
