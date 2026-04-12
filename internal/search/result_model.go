package search

import (
	"fmt"
	"slices"

	common "github.com/skylissh/melodihub/internal/domain"
)

type ResultInfo interface {
	GetID() string
	GetName() string
	GetType() common.InfoType
	GetImages() []common.Image
	GetMatch() int
}

type Result struct {
	TopResult *TopResult      `json:"top_result"`
	Tracks    []TrackSummary  `json:"tracks"`
	Artists   []ArtistSummary `json:"artists"`
	Albums    []AlbumSummary  `json:"albums"`
}

func (r *Result) FindTopResult() (*TopResult, error) {
	if r.TopResult.ID != "" {
		return r.TopResult, nil
	}

	length := len(r.Tracks) + len(r.Artists) + len(r.Albums)
	if length == 0 {
		return nil, fmt.Errorf("no results found")
	}
	results := make([]ResultInfo, 0, length)

	for _, track := range r.Tracks {
		results = append(results, track)
	}
	for _, artist := range r.Artists {
		results = append(results, artist)
	}
	for _, album := range r.Albums {
		results = append(results, album)
	}

	top := slices.MaxFunc(results, func(x, y ResultInfo) int {
		return x.GetMatch() - y.GetMatch()
	})

	return &TopResult{
		ID:     top.GetID(),
		Name:   top.GetName(),
		Type:   top.GetType(),
		Images: top.GetImages(),
		Match:  top.GetMatch(),
	}, nil
}
