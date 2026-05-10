package search

import (
	"fmt"

	"github.com/skylissh/melodihub/internal/domain"
	"github.com/skylissh/melodihub/internal/providers/deezer"
	"github.com/skylissh/melodihub/internal/utils"
)

type TrackSummaryResponse struct {
	ID       string                  `json:"id" validate:"required"`
	Type     domain.InfoType         `json:"type" validate:"required"`
	Duration int                     `json:"duration"`
	Title    string                  `json:"name" validate:"required"`
	Images   []domain.Image          `json:"images"`
	Artists  []ArtistSummaryResponse `json:"artists"`
	Match    int                     `json:"match"`
}

func TrackResponseFromSummary(summary TrackSummary) (TrackSummaryResponse, error) {
	var zero TrackSummaryResponse

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

	return TrackSummaryResponse{
		ID:       summary.id.String(),
		Type:     summary.kind,
		Duration: summary.duration,
		Title:    summary.title,
		Images:   summary.images,
		Artists:  artists,
		Match:    *summary.match,
	}, nil
}

func TrackResponseFromSlice(summaries []TrackSummary) ([]TrackSummaryResponse, error) {
	var responses []TrackSummaryResponse
	for _, summary := range summaries {
		response, err := TrackResponseFromSummary(summary)
		if err != nil {
			return nil, err
		}
		responses = append(responses, response)
	}
	return responses, nil
}

type TrackSummary struct {
	id       domain.ID
	kind     domain.InfoType
	duration int
	title    string
	images   []domain.Image
	artists  []ArtistSummary
	match    *int
}

func TrackSummaryFromDeezer(source deezer.Track) (TrackSummary, error) {
	id, err := domain.NewDeezerID(source.ID)
	if err != nil {
		return TrackSummary{}, err
	}

	return TrackSummary{
		id:       id,
		kind:     domain.TrackType,
		duration: source.Duration,
		title:    source.Title,
		images:   []domain.Image{{URL: source.Album.Cover}},
	}, nil
}

func (t *TrackSummary) SetArtists(artists []ArtistSummary) {
	t.artists = artists
}

// Matchable impl
func (t *TrackSummary) MatchByQuery(query Query) {
	match := utils.RankString(query.String(), fmt.Sprintf(
		"%s - %s",
		t.title, t.artists[0].name,
	))
	t.match = &match
}

// SearchItem impl
func (t TrackSummary) ID() domain.ID { return t.id }

func (t TrackSummary) Name() string { return t.title }

func (t TrackSummary) Kind() domain.InfoType { return t.kind }

func (t TrackSummary) Duration() int { return t.duration }

func (t TrackSummary) Match() *int { return t.match }

func (t TrackSummary) Images() []domain.Image { return t.images }

func (t TrackSummary) Artists() []ArtistSummary { return t.artists }
