package lastfm

type ArtistSimple struct {
	Name string `json:"name"`
	Mbid string `json:"mbid,omitempty"`
	URL  string `json:"url"`
}

type Track struct {
	Name      string       `json:"name" validate:"required"`
	Duration  int          `json:"duration,string,omitempty"`
	Mbid      string       `json:"mbid"`
	URL       string       `json:"url"`
	Artist    ArtistSimple `json:"artist"`
	Listeners int          `json:"listeners,string"`
}

type TopTracskAttr struct {
	Country    string `json:"country"`
	Page       int    `json:"page,string"`
	PerPage    int    `json:"perPage,string"`
	TotalPages int    `json:"totalPages,string"`
	Total      int    `json:"total,string"`
}

type TopTracks struct {
	Tracks   []Track       `json:"track" validate:"required,dive"`
	PageInfo TopTracskAttr `json:"@attr"`
}
