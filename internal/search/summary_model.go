package search

import common "github.com/skylissh/melodihub/internal/domain"

type ArtistSummary struct {
	ID     string          `json:"id" validate:"required"`
	Type   common.InfoType `json:"type" validate:"required"`
	Name   string          `json:"name" validate:"required"`
	Images []common.Image  `json:"images"`
	Match  int             `json:"match,omitempty"`
}

func (a ArtistSummary) GetID() string { return a.ID }

func (a ArtistSummary) GetName() string { return a.Name }

func (a ArtistSummary) GetType() common.InfoType { return a.Type }

func (a ArtistSummary) GetImages() []common.Image { return a.Images }

func (a ArtistSummary) GetMatch() int { return a.Match }

type TrackSummary struct {
	ID       string          `json:"id" validate:"required"`
	Type     common.InfoType `json:"type" validate:"required"`
	Duration int             `json:"duration"`
	Title    string          `json:"name" validate:"required"`
	Images   []common.Image  `json:"images"`
	Artists  []ArtistSummary `json:"artists"`
	Match    int             `json:"match,omitempty"`
}

func (t TrackSummary) GetID() string { return t.ID }

func (t TrackSummary) GetName() string { return t.Title }

func (t TrackSummary) GetType() common.InfoType { return t.Type }

func (t TrackSummary) GetImages() []common.Image { return t.Images }

func (t TrackSummary) GetMatch() int { return t.Match }

type AlbumSummary struct {
	ID      string          `json:"id" validate:"required"`
	Type    common.InfoType `json:"type" validate:"required"`
	Title   string          `json:"name" validate:"required"`
	Images  []common.Image  `json:"images"`
	Artists []ArtistSummary `json:"artists"`
	Match   int             `json:"match,omitempty"`
}

func (a AlbumSummary) GetID() string { return a.ID }

func (a AlbumSummary) GetName() string { return a.Title }

func (a AlbumSummary) GetType() common.InfoType { return a.Type }

func (a AlbumSummary) GetImages() []common.Image { return a.Images }

func (a AlbumSummary) GetMatch() int { return a.Match }
