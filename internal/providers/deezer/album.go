package deezer

type Album struct {
	ID    int    `json:"id" validate:"required"`
	Title string `json:"title" validate:"required"`
	Cover string `json:"cover_medium" validate:"required"`
}

type AlbumDetail struct {
	Album

	Genres         Data[Genre] `json:"genres"`
	Duration       int         `json:"duration"`
	Fans           int         `json:"fans"`
	ReleaseDate    string      `json:"release_date"`
	NumTracks      int         `json:"nb_tracks"`
	RecordType     string      `json:"record_type"`
	ExplicitLyrics bool        `json:"explicit_lyrics"`
	Artist         Artist      `json:"artist"`
	Tracks         Data[Track] `json:"tracks"`
}
