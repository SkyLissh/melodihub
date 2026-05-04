package deezer

type Track struct {
	ID             int    `json:"id" validate:"required"`
	Title          string `json:"title" validate:"required"`
	Link           string `json:"link" validate:"required"`
	Duration       int    `json:"duration"`
	Rank           int    `json:"rank"`
	ExplicitLyrics bool   `json:"explicit_lyrics"`
	Preview        string `json:"preview"`
	Album          Album  `json:"album"`
	Artist         Artist `json:"artist"`
}
