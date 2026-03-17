package model

type SimpleTrack struct {
	ID        string       `json:"id"`
	Duration  int          `json:"duration"`
	Title     string       `json:"title"`
	Artist    SimpleArtist `json:"artist"`
	Images    []Image      `json:"images,omitempty"`
	Rank      int          `json:"rank"`
	Listeners int          `json:"listeners,omitempty"`
}
