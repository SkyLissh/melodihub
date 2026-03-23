package deezer

type SearchResult struct {
	ID             int    `json:"id" validate:"required"`
	Readable       bool   `json:"readable" validate:"required"`
	Title          string `json:"title" validate:"required"`
	TitleShort     string `json:"title_short" validate:"required"`
	TitleVersion   string `json:"title_version" validate:"required"`
	Link           string `json:"link" validate:"required"`
	Duration       int    `json:"duration"`
	Rank           int    `json:"rank"`
	ExplicitLyrics bool   `json:"explicit_lyrics"`
	Preview        string `json:"preview"`
	Album          Album  `json:"album"`
	Artist         Artist `json:"artist"`
}
