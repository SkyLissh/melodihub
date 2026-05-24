package deezer

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound          = errors.New("not found")
	ErrItemsLimit        = errors.New("items limit exceeded")
	ErrPermission        = errors.New("permission denied")
	ErrTokenInvalid      = errors.New("token invalid")
	ErrParameterMissing  = errors.New("parameter missing")
	ErrInvalidParameter  = errors.New("invalid parameter")
	ErrQuotaExceeded     = errors.New("quota exceeded")
	ErrQueryInvalid      = errors.New("query invalid")
	ErrServiceBusy       = errors.New("service busy")
	ErrAccountNotAllowed = errors.New("individual account not allowed")
	ErrInvalidResponse   = errors.New("invalid response")
	ErrServer            = errors.New("server error")
	ErrUnknownCode       = errors.New("unknown error code")
)

type DeezerError struct {
	Kind error
	Msg  string
}

func (e DeezerError) Error() string {
	if e.Msg == "" {
		return e.Kind.Error()
	}

	return fmt.Sprintf("%v: %s", e.Kind, e.Msg)
}

func (e DeezerError) Unwrap() error {
	return e.Kind
}

func ServerError(err error) error {
	return DeezerError{
		Kind: ErrServer,
		Msg:  err.Error(),
	}
}

func InvalidResponse(err error) error {
	return DeezerError{
		Kind: ErrInvalidResponse,
		Msg:  err.Error(),
	}
}

func UnknownCode(code ErrorCode, kind string, msg string) error {
	return DeezerError{
		Kind: ErrUnknownCode,
		Msg:  fmt.Sprintf("code=%d kind=%s msg=%s", code, kind, msg),
	}
}
