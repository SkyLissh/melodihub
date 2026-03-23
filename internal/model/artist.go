package model

type SimpleArtist struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Images []Image `json:"images,omitempty"`
}

type Artist struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Listeners int            `json:"listeners"`
	Bio       string         `json:"bio,omitempty"`
	Images    []Image        `json:"images,omitempty"`
	TopTracks []SimpleTrack  `json:"top_tracks,omitempty"`
	TopAlbums []SimpleAlbum  `json:"albums,omitempty"`
	Similar   []SimpleArtist `json:"similar_artists,omitempty"`
}
