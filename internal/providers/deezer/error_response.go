package deezer

type ErrorResponse struct {
	Kind    string    `json:"type"`
	Message string    `json:"message"`
	Code    ErrorCode `json:"code"`
}

func (e ErrorResponse) ToError() error {
	switch e.Code {
	case QuoteExceededCode:
		return DeezerError{Kind: ErrQuotaExceeded, Msg: e.Message}
	case ItemsLimitExceededCode:
		return DeezerError{Kind: ErrItemsLimit, Msg: e.Message}
	case PermissionCode:
		return DeezerError{Kind: ErrPermission, Msg: e.Message}
	case TokenInvalidCode:
		return DeezerError{Kind: ErrTokenInvalid, Msg: e.Message}
	case ParameterInvalidCode:
		return DeezerError{Kind: ErrInvalidParameter, Msg: e.Message}
	case ParameterMissingCode:
		return DeezerError{Kind: ErrParameterMissing, Msg: e.Message}
	case QueryInvalidCode:
		return DeezerError{Kind: ErrQueryInvalid, Msg: e.Message}
	case ServiceBusyCode:
		return DeezerError{Kind: ErrServiceBusy, Msg: e.Message}
	case DataNotFoundCode:
		return DeezerError{Kind: ErrNotFound, Msg: e.Message}
	case AccountNotAllowedCode:
		return DeezerError{Kind: ErrAccountNotAllowed, Msg: e.Message}
	default:
		return UnknownCode(e.Code, e.Kind, e.Message)
	}
}
