package search

import (
	common "github.com/skylissh/melodihub/internal/domain"
)

type TopResult struct {
	ID      string          `json:"id"`
	Type    common.InfoType `json:"type"`
	Name    string          `json:"name"`
	Images  []common.Image  `json:"images,omitempty"`
	Match   int             `json:"match,omitempty"`
	Artists []ArtistSummary `json:"artists,omitempty"`
}

func (r TopResult) GetID() string             { return r.ID }
func (r TopResult) GetName() string           { return r.Name }
func (r TopResult) GetType() common.InfoType  { return r.Type }
func (r TopResult) GetImages() []common.Image { return r.Images }
func (r TopResult) GetMatch() int             { return r.Match }
