package deezer

type Response[T any] struct {
	Data T `json:"data" validate:"required"`
}
