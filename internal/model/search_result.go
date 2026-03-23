package model

type TrackResult struct {
	ID             string       `json:"id"`
	Title          string       `json:"title"`
	Duration       int          `json:"duration"`
	ExplicitLyrics bool         `json:"explicit_lyrics"`
	Preview        string       `json:"preview"`
	Rank           int          `json:"rank"`
	Artist         ArtistResult `json:"artist"`
	Match          int          `json:"match,omitempty"`
	Type           InfoType     `json:"type"`
	Image          string       `json:"image"`
}

type ArtistResult struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Image string   `json:"image"`
	Match int      `json:"match,omitempty"`
	Type  InfoType `json:"type"`
}

type AlbumResult struct {
	ID     string       `json:"id"`
	Title  string       `json:"title"`
	Image  string       `json:"image"`
	Artist ArtistResult `json:"artist"`
	Match  int          `json:"match,omitempty"`
	Type   InfoType     `json:"type"`
}

type SearchResult struct {
	Tracks  []TrackResult  `json:"tracks"`
	Artists []ArtistResult `json:"artists"`
	Albums  []AlbumResult  `json:"albums"`
}
