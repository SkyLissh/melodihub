package deezer

type Data[T any] struct {
	Data T `json:"data" validate:"required"`
}

type Response[T any] struct {
	Data  *T             `json:"-"`
	Error *ErrorResponse `json:"error,omitempty"`
}
