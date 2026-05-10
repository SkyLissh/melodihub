package lastfm

type ErrorResponse struct {
	Code    ErrorCode `json:"error"`
	Message string    `json:"message"`
}

func (e ErrorResponse) HasError() bool {
	return e.Code != 0
}

func (e ErrorResponse) ToError() error {
	switch e.Code {
	case InvalidServiceCode:
		return NewLastfmError(ErrInvalidService, e.Message)
	case InvalidMethodCode:
		return NewLastfmError(ErrInvalidMethod, e.Message)
	case AuthenticationFailedCode:
		return NewLastfmError(ErrAuthFailed, e.Message)
	case InvalidFormatCode:
		return NewLastfmError(ErrInvalidFormat, e.Message)
	case BadAuthTokenCode:
		return NewLastfmError(ErrBadAuthToken, e.Message)
	case InvalidParametersCode:
		return NewLastfmError(ErrInvalidParameter, e.Message)
	case InvalidResourceCode:
		return NewLastfmError(ErrInvalidResource, e.Message)
	case OperationFailedCode:
		return NewLastfmError(ErrOperationFailed, e.Message)
	case InvalidSessionKeyCode:
		return NewLastfmError(ErrInvalidSessionKey, e.Message)
	case InvalidMethodSignatureCode:
		return NewLastfmError(ErrInvalidSignature, e.Message)
	case TemporaryServerIssuesCode:
		return NewLastfmError(ErrTemporaryServerIssues, e.Message)
	case InvalidUsernameCode:
		return NewLastfmError(ErrInvalidUsername, e.Message)
	case InvalidTimestampCode:
		return NewLastfmError(ErrInvalidTimestamp, e.Message)
	case DeletedAPIKeyCode:
		return NewLastfmError(ErrDeletedAPIKey, e.Message)
	case NotEnoughContentCode:
		return NewLastfmError(ErrNotEnoughContent, e.Message)
	case ServiceUnavailableCode:
		return NewLastfmError(ErrServiceUnavailable, e.Message)
	case UserRequiredCode:
		return NewLastfmError(ErrUserRequired, e.Message)
	case SubsonicServerRequiredCode:
		return NewLastfmError(ErrSubsonicServerRequired, e.Message)
	case LegacyMethodDisabledCode:
		return NewLastfmError(ErrLegacyMethodDisabled, e.Message)
	case BadAccountScrobbleCode:
		return NewLastfmError(ErrBadAccountScrobble, e.Message)
	case NonExistentErrorCode:
		return NewLastfmError(ErrNonExistentError, e.Message)
	case RegistrationDisabledCode:
		return NewLastfmError(ErrRegistrationDisabled, e.Message)
	case TooManyRequestsCode:
		return NewLastfmError(ErrTooManyRequests, e.Message)
	case APIKeyPermissionDeniedCode:
		return NewLastfmError(ErrAPIKeyPermissionDenied, e.Message)
	case SuspendedAPIKeyCode:
		return NewLastfmError(ErrSuspendedAPIKey, e.Message)
	case RateLimitedCode:
		return NewLastfmError(ErrRateLimited, e.Message)
	case OfflineCode:
		return NewLastfmError(ErrOffline, e.Message)
	default:
		return UnknownCode(e.Code, e.Message)
	}
}
