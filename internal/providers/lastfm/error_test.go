package lastfm

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorResponseToErrorMapsEveryKnownCodeToSemanticError(t *testing.T) {
	tests := []struct {
		name string
		code ErrorCode
		err  error
	}{
		{name: "invalid service", code: InvalidServiceCode, err: ErrInvalidService},
		{name: "invalid method", code: InvalidMethodCode, err: ErrInvalidMethod},
		{name: "authentication failed", code: AuthenticationFailedCode, err: ErrAuthFailed},
		{name: "invalid format", code: InvalidFormatCode, err: ErrInvalidFormat},
		{name: "bad auth token", code: BadAuthTokenCode, err: ErrBadAuthToken},
		{name: "invalid parameters", code: InvalidParametersCode, err: ErrInvalidParameter},
		{name: "invalid resource", code: InvalidResourceCode, err: ErrInvalidResource},
		{name: "operation failed", code: OperationFailedCode, err: ErrOperationFailed},
		{name: "invalid session key", code: InvalidSessionKeyCode, err: ErrInvalidSessionKey},
		{name: "invalid method signature", code: InvalidMethodSignatureCode, err: ErrInvalidSignature},
		{name: "temporary server issues", code: TemporaryServerIssuesCode, err: ErrTemporaryServerIssues},
		{name: "invalid username", code: InvalidUsernameCode, err: ErrInvalidUsername},
		{name: "invalid timestamp", code: InvalidTimestampCode, err: ErrInvalidTimestamp},
		{name: "deleted api key", code: DeletedAPIKeyCode, err: ErrDeletedAPIKey},
		{name: "not enough content", code: NotEnoughContentCode, err: ErrNotEnoughContent},
		{name: "service unavailable", code: ServiceUnavailableCode, err: ErrServiceUnavailable},
		{name: "user required", code: UserRequiredCode, err: ErrUserRequired},
		{name: "subsonic server required", code: SubsonicServerRequiredCode, err: ErrSubsonicServerRequired},
		{name: "legacy method disabled", code: LegacyMethodDisabledCode, err: ErrLegacyMethodDisabled},
		{name: "bad account scrobble", code: BadAccountScrobbleCode, err: ErrBadAccountScrobble},
		{name: "non existent error", code: NonExistentErrorCode, err: ErrNonExistentError},
		{name: "registration disabled", code: RegistrationDisabledCode, err: ErrRegistrationDisabled},
		{name: "too many requests", code: TooManyRequestsCode, err: ErrTooManyRequests},
		{name: "api key permission denied", code: APIKeyPermissionDeniedCode, err: ErrAPIKeyPermissionDenied},
		{name: "suspended api key", code: SuspendedAPIKeyCode, err: ErrSuspendedAPIKey},
		{name: "rate limited", code: RateLimitedCode, err: ErrRateLimited},
		{name: "offline", code: OfflineCode, err: ErrOffline},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrorResponse{
				Code:    tt.code,
				Message: "provider message",
			}.ToError()

			assert.ErrorIs(t, err, tt.err)
			assert.Contains(t, err.Error(), "provider message")
		})
	}
}

func TestErrorResponseToErrorPreservesUnknownCode(t *testing.T) {
	err := ErrorResponse{
		Code:    ErrorCode(999),
		Message: "new Last.fm error",
	}.ToError()

	require.Error(t, err)
	assert.ErrorIs(t, err, ErrUnknownCode)
	assert.Contains(t, err.Error(), "code=999")
	assert.Contains(t, err.Error(), "new Last.fm error")
}

func TestLastfmErrorUnwrapSupportsErrorsIs(t *testing.T) {
	err := NewLastfmError(ErrRateLimited, "Rate Limit Exceeded")

	assert.True(t, errors.Is(err, ErrRateLimited))
	assert.Equal(t, "rate limited: Rate Limit Exceeded", err.Error())
}
