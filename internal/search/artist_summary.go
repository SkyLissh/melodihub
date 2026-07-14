package search

import (
	"fmt"

	"github.com/skylissh/melodihub/internal/domain"
	"github.com/skylissh/melodihub/internal/providers/deezer"
	"github.com/skylissh/melodihub/internal/utils"
)

type ArtistSummaryResponse struct {
	ID     string          `json:"id" validate:"required"`
	Kind   domain.InfoType `json:"kind" validate:"required"`
	Name   string          `json:"name" validate:"required"`
	Images []domain.Image  `json:"images"`
	Match  int             `json:"match"`
}

func ArtistResponseFromSummary(summary ArtistSummary) (ArtistSummaryResponse, error) {
	var zero ArtistSummaryResponse

	if summary.match == nil {
		return zero, fmt.Errorf("match cannot be nil")
	}

	return ArtistSummaryResponse{
		ID:     summary.id.String(),
		Kind:   summary.kind,
		Name:   summary.name,
		Images: summary.images,
		Match:  *summary.match,
	}, nil
}

func ArtistResponseFromSlice(summaries []ArtistSummary) ([]ArtistSummaryResponse, error) {
	var responses []ArtistSummaryResponse
	for _, summary := range summaries {
		response, err := ArtistResponseFromSummary(summary)
		if err != nil {
			return nil, err
		}
		responses = append(responses, response)
	}
	return responses, nil
}

type ArtistSummary struct {
	id     domain.ID
	name   string
	kind   domain.InfoType
	images []domain.Image
	match  *int
}

func ArtistSummaryFromDeezer(source deezer.Artist) (ArtistSummary, error) {
	id, err := domain.NewDeezerID(source.ID)
	if err != nil {
		return ArtistSummary{}, err
	}

	return ArtistSummary{
		id:     id,
		name:   source.Name,
		kind:   domain.ArtistType,
		images: []domain.Image{{URL: source.Picture}},
	}, nil
}

// Matchable impl
func (a *ArtistSummary) MatchByQuery(query Query) {
	match := utils.RankString(query.String(), a.name)
	a.match = &match
}

// SearchItem impl
func (a ArtistSummary) ID() domain.ID { return a.id }

func (a ArtistSummary) Name() string { return a.name }

func (a ArtistSummary) Kind() domain.InfoType { return a.kind }

func (a ArtistSummary) Images() []domain.Image { return a.images }

func (a ArtistSummary) Match() *int { return a.match }

func (a ArtistSummary) Artists() []ArtistSummary { return []ArtistSummary{a} }
