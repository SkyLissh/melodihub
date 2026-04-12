package carousel

import (
	common "github.com/skylissh/melodihub/internal/domain"
)

type ArtistSummary struct {
	ID     string          `json:"id" validate:"required"`
	Type   common.InfoType `json:"type" validate:"required"`
	Name   string          `json:"name" validate:"required"`
	Images []common.Image  `json:"images"`
}

type Item struct {
	Type    common.InfoType `json:"type"`
	ID      string          `json:"id"`
	Name    string          `json:"name"`
	Images  []common.Image  `json:"images,omitempty"`
	Artists []ArtistSummary `json:"artists,omitempty"`
}
