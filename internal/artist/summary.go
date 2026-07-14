package artist

import (
	"github.com/skylissh/melodihub/internal/domain"
	"github.com/skylissh/melodihub/internal/providers/deezer"
)

type SummaryResponse struct {
	ID     string          `json:"id" validate:"required"`
	Kind   domain.InfoType `json:"kind" validate:"required"`
	Name   string          `json:"name" validate:"required"`
	Images []domain.Image  `json:"images"`
}

func ResponseFromSummary(summary *Summary) SummaryResponse {
	return SummaryResponse{
		ID:     summary.ID.String(),
		Kind:   summary.Kind,
		Name:   summary.Name,
		Images: summary.Images,
	}
}

func ResponseFromSummaries(summaries []Summary) []SummaryResponse {
	var result []SummaryResponse
	for _, summary := range summaries {
		result = append(result, ResponseFromSummary(&summary))
	}
	return result
}

type Summary struct {
	ID     domain.ID
	Kind   domain.InfoType
	Name   string
	Images []domain.Image
}

func SummaryFromDeezer(source deezer.Artist) (Summary, error) {
	id, err := domain.NewDeezerID(source.ID)
	if err != nil {
		return Summary{}, err
	}
	return Summary{
		ID:     id,
		Kind:   domain.ArtistType,
		Name:   source.Name,
		Images: []domain.Image{{URL: source.Picture}},
	}, nil
}
