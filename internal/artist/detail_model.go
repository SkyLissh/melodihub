package artist

import (
	common "github.com/skylissh/melodihub/internal/domain"
)

type Detail struct {
	Summary

	Listeners int        `json:"listeners" validate:"required,min=10"`
	Bio       string     `json:"bio,omitempty"`
	TopTracks []TopTrack `json:"top_tracks"`
	TopAlbums []TopAlbum `json:"top_albums"`
	Similar   []Summary  `json:"similar_artists"`
}

type TopTrack struct {
	ID        string          `json:"id" validate:"required"`
	Type      common.InfoType `json:"type" validate:"required"`
	Title     string          `json:"title" validate:"required"`
	Images    []common.Image  `json:"images"`
	Duration  int             `json:"duration" validate:"min=1"`
	Listeners int             `json:"listeners" validate:"min=1"`
}

type TopAlbum struct {
	ID     string          `json:"id" validate:"required"`
	Type   common.InfoType `json:"type" validate:"required"`
	Title  string          `json:"title" validate:"required"`
	Images []common.Image  `json:"images"`
}
