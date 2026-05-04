package deezer

type Genre struct {
	ID      int    `json:"id" validate:"required"`
	Name    string `json:"name" validate:"required"`
	Picture string `json:"picture_medium" validate:"required"`
}
