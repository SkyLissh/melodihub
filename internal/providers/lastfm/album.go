package lastfm

type Album struct {
	Mbid      string       `json:"mbid,omitempty"`
	Name      string       `json:"name" validate:"required"`
	Playcount int          `json:"playcount"`
	Artist    SimpleArtist `json:"artist"`
}

type TopAlbumAttr struct {
	Page       int `json:"page,string"`
	PerPage    int `json:"perPage,string"`
	TotalPages int `json:"totalPages,string"`
	Total      int `json:"total,string"`
}

type TopAlbum struct {
	Album []Album      `json:"album" validate:"required,dive"`
	Attr  TopAlbumAttr `json:"@attr"`
}
