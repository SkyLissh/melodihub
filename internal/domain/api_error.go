package domain

import (
	"errors"
	"net/http"
)

type APIError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Cause   string    `json:"cause,omitempty"`
}

func NewAPIError(message string, code int) APIError {
	errorCode := CodeInternal
	if code >= http.StatusBadRequest && code < http.StatusInternalServerError {
		errorCode = CodeInvalidParam
	}

	return APIError{
		Code:    errorCode,
		Message: message,
	}
}

func APIErrorFromMelodiError(err error) (APIError, int) {
	var melodiErr *MelodiError
	if !errors.As(err, &melodiErr) {
		return APIError{
			Code:    CodeInternal,
			Message: "Internal error",
		}, http.StatusInternalServerError
	}

	apiErr := APIError{
		Code:    melodiErr.Code,
		Message: melodiErr.Message,
		Cause:   melodiErr.Cause,
	}

	return apiErr, statusFromErrorCode(melodiErr.Code)
}

func statusFromErrorCode(code ErrorCode) int {
	switch code {
	case CodeInvalidParam, CodeInvalidBody:
		return http.StatusBadRequest
	case CodeNotFound:
		return http.StatusNotFound
	case CodeRateLimited:
		return http.StatusTooManyRequests
	case CodeProviderUnavailable:
		return http.StatusBadGateway
	default:
		return http.StatusInternalServerError
	}
}
