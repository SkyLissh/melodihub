package domain

type APIError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func NewAPIError(message string, code int) APIError {
	return APIError{
		Message: message,
		Code:    code,
	}
}
