package lastfm

import "fmt"

type ErrorCode int

const (
	InvalidServiceCode         ErrorCode = 1
	InvalidMethodCode          ErrorCode = 2
	AuthenticationFailedCode   ErrorCode = 3
	InvalidFormatCode          ErrorCode = 4
	BadAuthTokenCode           ErrorCode = 5
	InvalidParametersCode      ErrorCode = 6
	InvalidResourceCode        ErrorCode = 7
	OperationFailedCode        ErrorCode = 8
	InvalidSessionKeyCode      ErrorCode = 9
	InvalidMethodSignatureCode ErrorCode = 10
	TemporaryServerIssuesCode  ErrorCode = 11
	InvalidUsernameCode        ErrorCode = 13
	InvalidTimestampCode       ErrorCode = 14
	DeletedAPIKeyCode          ErrorCode = 15
	NotEnoughContentCode       ErrorCode = 16
	ServiceUnavailableCode     ErrorCode = 17
	UserRequiredCode           ErrorCode = 18
	SubsonicServerRequiredCode ErrorCode = 19
	LegacyMethodDisabledCode   ErrorCode = 20
	BadAccountScrobbleCode     ErrorCode = 21
	NonExistentErrorCode       ErrorCode = 22
	RegistrationDisabledCode   ErrorCode = 23
	TooManyRequestsCode        ErrorCode = 24
	APIKeyPermissionDeniedCode ErrorCode = 25
	SuspendedAPIKeyCode        ErrorCode = 26
	RateLimitedCode            ErrorCode = 27
	OfflineCode                ErrorCode = 29
)

var errorCodeNames = map[ErrorCode]string{
	InvalidServiceCode:         "InvalidService",
	InvalidMethodCode:          "InvalidMethod",
	AuthenticationFailedCode:   "AuthenticationFailed",
	InvalidFormatCode:          "InvalidFormat",
	BadAuthTokenCode:           "BadAuthToken",
	InvalidParametersCode:      "InvalidParameters",
	InvalidResourceCode:        "InvalidResource",
	OperationFailedCode:        "OperationFailed",
	InvalidSessionKeyCode:      "InvalidSessionKey",
	InvalidMethodSignatureCode: "InvalidMethodSignature",
	TemporaryServerIssuesCode:  "TemporaryServerIssues",
	InvalidUsernameCode:        "InvalidUsername",
	InvalidTimestampCode:       "InvalidTimestamp",
	DeletedAPIKeyCode:          "DeletedAPIKey",
	NotEnoughContentCode:       "NotEnoughContent",
	ServiceUnavailableCode:     "ServiceUnavailable",
	UserRequiredCode:           "UserRequired",
	SubsonicServerRequiredCode: "SubsonicServerRequired",
	LegacyMethodDisabledCode:   "LegacyMethodDisabled",
	BadAccountScrobbleCode:     "BadAccountScrobble",
	NonExistentErrorCode:       "NonExistentError",
	RegistrationDisabledCode:   "RegistrationDisabled",
	TooManyRequestsCode:        "TooManyRequests",
	APIKeyPermissionDeniedCode: "APIKeyPermissionDenied",
	SuspendedAPIKeyCode:        "SuspendedAPIKey",
	RateLimitedCode:            "RateLimited",
	OfflineCode:                "Offline",
}

func (e ErrorCode) String() string {
	if name, ok := errorCodeNames[e]; ok {
		return name
	}

	return fmt.Sprintf("Unknown(%d)", e)
}

func (e ErrorCode) IsKnown() bool {
	_, ok := errorCodeNames[e]
	return ok
}
