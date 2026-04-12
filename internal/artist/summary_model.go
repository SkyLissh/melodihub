package artist

import common "github.com/skylissh/melodihub/internal/domain"

type Summary struct {
	ID     string          `json:"id" validate:"required"`
	Type   common.InfoType `json:"type" validate:"required"`
	Name   string          `json:"name" validate:"required"`
	Images []common.Image  `json:"images"`
}

func NewSummary(id, name string, images []common.Image) Summary {
	return Summary{
		ID:     id,
		Type:   common.ArtistType,
		Name:   name,
		Images: images,
	}
}
