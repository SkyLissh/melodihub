package album

import common "github.com/skylissh/melodihub/internal/domain"

type Detail struct {
	ID          string          `json:"id"`
	Title       string          `json:"title"`
	Tracks      []Track         `json:"tracks"`
	Artist      ArtistSummary   `json:"artist"`
	Images      []common.Image  `json:"images"`
	Features    []ArtistSummary `json:"features"`
	ReleaseDate string          `json:"release_date"`
	Duration    int             `json:"duration"`
}

type Track struct {
	ID        string          `json:"id"`
	Title     string          `json:"title"`
	Duration  int             `json:"duration"`
	Artists   []ArtistSummary `json:"artists"`
	PlayCount int             `json:"play_count"`
}

type ArtistSummary struct {
	ID     string         `json:"id"`
	Name   string         `json:"name"`
	Images []common.Image `json:"images"`
}
