package search

import (
	"fmt"

	"github.com/skylissh/melodihub/internal/domain"
)

type TopResultResponse struct {
	ID      string                  `json:"id" validate:"required"`
	Kind    domain.InfoType         `json:"kind" validate:"required"`
	Name    string                  `json:"name" validate:"required"`
	Images  []domain.Image          `json:"images,omitempty"`
	Match   int                     `json:"match,omitempty"`
	Artists []ArtistSummaryResponse `json:"artists,omitempty"`
}

func TopResultResponseFromSummary(summary TopResult) (TopResultResponse, error) {
	var zero TopResultResponse

	if summary.match == nil {
		return zero, fmt.Errorf("match cannot be nil")
	}

	if summary.artists == nil {
		return zero, fmt.Errorf("artists cannot be nil")
	}

	if len(summary.artists) == 0 {
		return zero, fmt.Errorf("artists cannot be empty")
	}

	artists, err := ArtistResponseFromSlice(summary.artists)
	if err != nil {
		return zero, err
	}

	return TopResultResponse{
		ID:      summary.id.String(),
		Kind:    summary.kind,
		Name:    summary.name,
		Images:  summary.images,
		Match:   *summary.match,
		Artists: artists,
	}, nil
}

type TopResult struct {
	id      domain.ID
	name    string
	kind    domain.InfoType
	images  []domain.Image
	artists []ArtistSummary
	match   *int
}

func TopResultFromSearchItem(source SearchItem) (TopResult, error) {
	match := source.Match()
	if match == nil {
		return TopResult{}, fmt.Errorf("match cannot be nil")
	}

	return TopResult{
		id:      source.ID(),
		name:    source.Name(),
		kind:    source.Kind(),
		images:  source.Images(),
		artists: source.Artists(),
		match:   match,
	}, nil
}
