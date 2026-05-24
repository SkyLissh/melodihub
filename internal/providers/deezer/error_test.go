package deezer

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
		{name: "quota exceeded", code: QuoteExceededCode, err: ErrQuotaExceeded},
		{name: "items limit exceeded", code: ItemsLimitExceededCode, err: ErrItemsLimit},
		{name: "permission denied", code: PermissionCode, err: ErrPermission},
		{name: "token invalid", code: TokenInvalidCode, err: ErrTokenInvalid},
		{name: "invalid parameter", code: ParameterInvalidCode, err: ErrInvalidParameter},
		{name: "parameter missing", code: ParameterMissingCode, err: ErrParameterMissing},
		{name: "query invalid", code: QueryInvalidCode, err: ErrQueryInvalid},
		{name: "service busy", code: ServiceBusyCode, err: ErrServiceBusy},
		{name: "not found", code: DataNotFoundCode, err: ErrNotFound},
		{name: "account not allowed", code: AccountNotAllowedCode, err: ErrAccountNotAllowed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrorResponse{
				Code:    tt.code,
				Kind:    tt.name,
				Message: "provider message",
			}.ToError()

			assert.ErrorIs(t, err, tt.err)

		var deezerErr DeezerError
		require.ErrorAs(t, err, &deezerErr)
		assert.Equal(t, tt.err, deezerErr.Kind)
		assert.Equal(t, "provider message", deezerErr.Msg)
		})
	}
}

func TestErrorResponseToErrorPreservesUnknownCode(t *testing.T) {
	err := ErrorResponse{
		Code:    ErrorCode(999),
		Kind:    "NewError",
		Message: "new Deezer error",
	}.ToError()

	require.Error(t, err)
	assert.ErrorIs(t, err, ErrUnknownCode)
	assert.Contains(t, err.Error(), "code=999")
	assert.Contains(t, err.Error(), "kind=NewError")
	assert.Contains(t, err.Error(), "new Deezer error")
}

func TestDeezerErrorUnwrapSupportsErrorsIs(t *testing.T) {
	err := DeezerError{
		Kind: ErrServer,
		Msg:  "connection refused",
	}

	assert.True(t, errors.Is(err, ErrServer))
	assert.Equal(t, "server error: connection refused", err.Error())
}
