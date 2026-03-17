package model

type SimpleAlbum struct {
	ID     string       `json:"id"`
	Title  string       `json:"title"`
	Artist SimpleArtist `json:"artist"`
	Images []Image      `json:"images,omitempty"`
}
