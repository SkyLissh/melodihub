package spotify

type SearchResult struct {
	Albums  *Paginated[AlbumSimple] `json:"albums,omitempty"`
	Artists *Paginated[Artist]      `json:"artists,omitempty"`
	Tracks  *Paginated[Track]       `json:"tracks,omitempty"`
}
