package lastfm

type ArtistStats struct {
	Listeners int `json:"listeners,string"`
	Playcount int `json:"playcount,string"`
}

type ArtistBio struct {
	Published string `json:"published"`
	Summary   string `json:"summary"`
	Content   string `json:"content"`
}

type ArtistDetail struct {
	Name       string      `json:"name"`
	Mbid       string      `json:"mbid,omitempty"`
	URL        string      `json:"url"`
	Streamable int         `json:"streamable,string"`
	OnTour     int         `json:"ontour,string"`
	Stats      ArtistStats `json:"stats"`
	Similar    struct {
		Artists []Artist `json:"artist"`
	} `json:"similar"`
	Tags struct {
		Tags []Tag `json:"tag"`
	} `json:"tags"`
	Bio ArtistBio `json:"bio"`
}

type SimpleArtist struct {
	Name string `json:"name"`
	Mbid string `json:"mbid,omitempty"`
	URL  string `json:"url"`
}

type Artist struct {
	Name       string `json:"name"`
	Listeners  int    `json:"listeners,string"`
	Mbid       string `json:"mbid,omitempty"`
	URL        string `json:"url"`
	Streamable string `json:"streamable"`
}

type TopArtistsAttr struct {
	Country    string `json:"country"`
	Page       int    `json:"page,string"`
	PerPage    int    `json:"perPage,string"`
	TotalPages int    `json:"totalPages,string"`
	Total      int    `json:"total,string"`
}

type ArtistsContainer struct {
	Artists  []Artist       `json:"artist"`
	PageInfo TopArtistsAttr `json:"@attr"`
}

type TopArtists struct {
	TopArtists ArtistsContainer `json:"topartists"`
}
