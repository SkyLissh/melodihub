package deezer

type Artist struct {
	ID      int    `json:"id" validate:"required"`
	Name    string `json:"name" validate:"required"`
	Link    string `json:"link,omitempty"`
	Picture string `json:"picture_medium" validate:"required"`
}

type TrackArtist struct {
	ID   int    `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}
