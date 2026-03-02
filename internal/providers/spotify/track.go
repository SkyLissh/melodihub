package spotify

type TrackSimple struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	URL         string         `json:"url"`
	Href        string         `json:"href"`
	TrackNumber int            `json:"track_number"`
	Explicit    bool           `json:"explicit"`
	DurationMs  int            `json:"duration_ms"`
	Artists     []ArtistSimple `json:"artists"`
	Type        InfoType       `json:"type"`
}

type Track struct {
	TrackSimple

	Popularity int         `json:"popularity"`
	Artists    []Artist    `json:"artists"`
	Album      AlbumSimple `json:"album"`
}
