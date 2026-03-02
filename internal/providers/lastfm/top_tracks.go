package lastfm

type ArtistSimple struct {
	Name string `json:"name"`
	Mbid string `json:"mbid,omitempty"`
	URL  string `json:"url"`
}

type Streamable struct {
	Text string `json:"#text"`
	Full int    `json:"fulltrack"`
}

type Track struct {
	Name       string       `json:"name"`
	Duration   int          `json:"duration"`
	Mbid       string       `json:"mbid"`
	URL        string       `json:"url"`
	Streamable Streamable   `json:"streamable"`
	Artist     ArtistSimple `json:"artist"`
	Listeners  int          `json:"listeners"`
}

type TopTracskAttr struct {
	Country    string `json:"country"`
	Page       int    `json:"page"`
	PerPage    int    `json:"perPage"`
	TotalPages int    `json:"totalPages"`
	Total      int    `json:"total"`
}

type TrackContainer struct {
	Tracks   []Track       `json:"track"`
	PageInfo TopTracskAttr `json:"@attr"`
}

type TopTracks struct {
	TopTracks TrackContainer `json:"toptracks"`
}
