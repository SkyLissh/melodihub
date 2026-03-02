package lastfm

type Artist struct {
	Name       string `json:"name"`
	Listeners  int    `json:"listeners"`
	Mbid       string `json:"mbid,omitempty"`
	URL        string `json:"url"`
	Streamable string `json:"streamable"`
}

type TopArtistsAttr struct {
	Country    string `json:"country"`
	Page       int    `json:"page"`
	PerPage    int    `json:"perPage"`
	TotalPages int    `json:"totalPages"`
	Total      int    `json:"total"`
}

type ArtistsContainer struct {
	Artists  []Artist       `json:"artist"`
	PageInfo TopArtistsAttr `json:"@attr"`
}

type TopArtists struct {
	TopArtists ArtistsContainer `json:"topartists"`
}
