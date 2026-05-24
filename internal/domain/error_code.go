package domain

type ErrorCode string

const (
	CodeInvalidParam        ErrorCode = "INVALID_PARAM"
	CodeInvalidBody         ErrorCode = "INVALID_BODY"
	CodeNotFound            ErrorCode = "NOT_FOUND"
	CodeProviderUnavailable ErrorCode = "PROVIDER_UNAVAILABLE"
	CodeRateLimited         ErrorCode = "RATE_LIMITED"
	CodeInternal            ErrorCode = "INTERNAL"
)
