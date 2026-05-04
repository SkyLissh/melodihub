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

func mapDeezerError(res *ErrorResponse) error {
	err := DeezerError{}

	switch res.Code {
	case 4:
		err.Kind = ErrQuotaExceeded
	case 100:
		err.Kind = ErrItemsLimit
	case 200:
		err.Kind = ErrPermission
	case 300:
		err.Kind = ErrTokenInvalid
	case 500:
		err.Kind = ErrInvalidParameter
	case 501:
		err.Kind = ErrParameterMissing
	case 600:
		err.Kind = ErrQueryInvalid
	case 700:
		err.Kind = ErrServiceBusy
	case 800:
		err.Kind = ErrNotFound
	case 901:
		err.Kind = ErrAccountNotAllowed
	default:
		err.Kind = ErrServer
		err.Msg = fmt.Sprintf("code=%d message=%s", res.Code, res.Message)
	}

	return err
}

func NotFound(kind Kind, id string) error {
	return DeezerError{
		Kind: ErrNotFound,
		Msg:  fmt.Sprintf("%s with id %q", kind, id),
	}
}

func ParameterMissing(name string) error {
	return DeezerError{
		Kind: ErrParameterMissing,
		Msg:  name,
	}
}

func InvalidParameter(name string) error {
	return DeezerError{
		Kind: ErrInvalidParameter,
		Msg:  name,
	}
}

func QueryInvalid(query string) error {
	return DeezerError{
		Kind: ErrQueryInvalid,
		Msg:  query,
	}
}

func InvalidResponse(err error) error {
	return DeezerError{
		Kind: ErrInvalidResponse,
		Msg:  err.Error(),
	}
}

func ServerError(msg string) error {
	return DeezerError{
		Kind: ErrServer,
		Msg:  msg,
	}
}
