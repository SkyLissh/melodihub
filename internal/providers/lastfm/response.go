package lastfm

import "encoding/json"

type Response[T any] struct {
	Data  T
	Error *ErrorResponse
}

func (r *Response[T]) UnmarshalJSON(data []byte) error {
	var apiErr ErrorResponse
	if err := json.Unmarshal(data, &apiErr); err == nil && apiErr.HasError() {
		r.Error = &apiErr
		return nil
	}

	var result T
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}

	r.Data = result
	return nil
}
