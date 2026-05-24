package domain

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvalidParamCreatesPublicCause(t *testing.T) {
	err := InvalidParam(errors.New("search query cannot be empty"))

	var melodiErr *MelodiError
	require.ErrorAs(t, err, &melodiErr)
	assert.Equal(t, CodeInvalidParam, melodiErr.Code)
	assert.Equal(t, "Invalid param", melodiErr.Message)
	assert.Equal(t, "search query cannot be empty", melodiErr.Cause)
}

func TestInternalDoesNotExposeCause(t *testing.T) {
	err := Internal(errors.New("database password leaked in stack"))

	var melodiErr *MelodiError
	require.ErrorAs(t, err, &melodiErr)
	assert.Equal(t, CodeInternal, melodiErr.Code)
	assert.Equal(t, "Internal error", melodiErr.Message)
	assert.Empty(t, melodiErr.Cause)
}

func TestAPIErrorFromMelodiErrorMapsKnownErrors(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		status int
		apiErr APIError
	}{
		{
			name:   "invalid param",
			err:    InvalidParam(errors.New("search query cannot be empty")),
			status: http.StatusBadRequest,
			apiErr: APIError{
				Code:    CodeInvalidParam,
				Message: "Invalid param",
				Cause:   "search query cannot be empty",
			},
		},
		{
			name:   "invalid body",
			err:    InvalidBody(errors.New("invalid json")),
			status: http.StatusBadRequest,
			apiErr: APIError{
				Code:    CodeInvalidBody,
				Message: "Invalid body",
				Cause:   "invalid json",
			},
		},
		{
			name:   "not found",
			err:    NotFound(errors.New("track missing")),
			status: http.StatusNotFound,
			apiErr: APIError{
				Code:    CodeNotFound,
				Message: "Not found",
				Cause:   "track missing",
			},
		},
		{
			name:   "provider unavailable",
			err:    ProviderUnavailable(ProviderDeezer, errors.New("service busy")),
			status: http.StatusBadGateway,
			apiErr: APIError{
				Code:    CodeProviderUnavailable,
				Message: "Provider unavailable",
				Cause:   "deezer: service busy",
			},
		},
		{
			name:   "rate limited",
			err:    RateLimited(ProviderDeezer, errors.New("quota exceeded")),
			status: http.StatusTooManyRequests,
			apiErr: APIError{
				Code:    CodeRateLimited,
				Message: "Rate limited",
				Cause:   "deezer: quota exceeded",
			},
		},
		{
			name:   "internal",
			err:    Internal(errors.New("private detail")),
			status: http.StatusInternalServerError,
			apiErr: APIError{
				Code:    CodeInternal,
				Message: "Internal error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiErr, status := APIErrorFromMelodiError(tt.err)

			assert.Equal(t, tt.status, status)
			assert.Equal(t, tt.apiErr, apiErr)
		})
	}
}

func TestAPIErrorFromMelodiErrorMapsUnknownErrorToInternal(t *testing.T) {
	apiErr, status := APIErrorFromMelodiError(errors.New("unexpected"))

	assert.Equal(t, http.StatusInternalServerError, status)
	assert.Equal(t, APIError{
		Code:    CodeInternal,
		Message: "Internal error",
	}, apiErr)
}
