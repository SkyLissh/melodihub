package deezer

import "encoding/json"

type Data[T any] struct {
	Data T `json:"data" validate:"required"`
}

type Response[T any] struct {
	Data  *T
	Error *ErrorResponse `json:"error,omitempty"`
}

func (r *Response[T]) UnmarshalJSON(data []byte) error {
	var apiErr struct {
		Error *ErrorResponse `json:"error,omitempty"`
	}
	if err := json.Unmarshal(data, &apiErr); err == nil && apiErr.Error != nil {
		r.Error = apiErr.Error
		return nil
	}

	var result T
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}

	r.Data = &result
	return nil
}
