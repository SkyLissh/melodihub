package domain

import "fmt"

type MelodiError struct {
	Code    ErrorCode
	Message string
	Cause   string
	Err     error
}

func (e *MelodiError) Error() string {
	if e.Err == nil {
		return e.Message
	}

	return fmt.Sprintf("%s: %v", e.Message, e.Err)
}

func (e *MelodiError) Unwrap() error {
	return e.Err
}

func InvalidParam(err error) error {
	return newMelodiError(CodeInvalidParam, "Invalid param", publicCause(err), err)
}

func InvalidBody(err error) error {
	return newMelodiError(CodeInvalidBody, "Invalid body", publicCause(err), err)
}

func NotFound(err error) error {
	return newMelodiError(CodeNotFound, "Not found", publicCause(err), err)
}

func ProviderUnavailable(provider Provider, err error) error {
	return newMelodiError(CodeProviderUnavailable, "Provider unavailable", providerCause(provider, err), err)
}

func RateLimited(provider Provider, err error) error {
	return newMelodiError(CodeRateLimited, "Rate limited", providerCause(provider, err), err)
}

func Internal(err error) error {
	return newMelodiError(CodeInternal, "Internal error", "", err)
}

func newMelodiError(code ErrorCode, message string, cause string, err error) error {
	return &MelodiError{
		Code:    code,
		Message: message,
		Cause:   cause,
		Err:     err,
	}
}

func publicCause(err error) string {
	if err == nil {
		return ""
	}

	return err.Error()
}

func providerCause(provider Provider, err error) string {
	if err == nil {
		return provider.String()
	}

	return fmt.Sprintf("%s: %s", provider, err.Error())
}
