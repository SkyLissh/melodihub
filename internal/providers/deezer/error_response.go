package deezer

type ErrorResponse struct {
	Kind    string    `json:"type"`
	Message string    `json:"message"`
	Code    ErrorCode `json:"code"`
}

func (e ErrorResponse) ToError() error {
	switch e.Code {
	case QuoteExceededCode:
		return ErrQuotaExceeded
	case ItemsLimitExceededCode:
		return ErrItemsLimit
	case PermissionCode:
		return ErrPermission
	case TokenInvalidCode:
		return ErrTokenInvalid
	case ParameterInvalidCode:
		return ErrInvalidParameter
	case ParameterMissingCode:
		return ErrParameterMissing
	case QueryInvalidCode:
		return ErrQueryInvalid
	case ServiceBusyCode:
		return ErrServiceBusy
	case DataNotFoundCode:
		return ErrNotFound
	case AccountNotAllowedCode:
		return ErrAccountNotAllowed
	default:
		return UnknownCode(e.Code, e.Kind, e.Message)
	}
}
