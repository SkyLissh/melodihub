package lastfm

type SimilarTrack struct {
	Name      string       `json:"name"`
	Duration  int          `json:"duration"`
	Mbid      string       `json:"mbid"`
	URL       string       `json:"url"`
	PlayCount int          `json:"playcount"`
	Artist    ArtistSimple `json:"artist"`
	Match     float64      `json:"match"`
}

type SimilaarTrackAttr struct {
	Track  string `json:"track"`
	Artist string `json:"artist"`
}

type SimilarTracks struct {
	Tracks []SimilarTrack `json:"track"`
	Attr   SimilaarTrackAttr
}
