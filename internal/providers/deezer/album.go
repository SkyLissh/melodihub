package deezer

type Album struct {
	ID    int    `json:"id" validate:"required"`
	Title string `json:"title" validate:"required"`
	Cover string `json:"cover_medium" validate:"required"`
}
