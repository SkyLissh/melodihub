package search

import (
	"fmt"

	"github.com/skylissh/melodihub/internal/domain"
	"github.com/skylissh/melodihub/internal/providers/deezer"
	"github.com/skylissh/melodihub/internal/utils"
)

type AlbumSummaryResponse struct {
	ID      string                  `json:"id" validate:"required"`
	Type    domain.InfoType         `json:"type" validate:"required"`
	Title   string                  `json:"name" validate:"required"`
	Images  []domain.Image          `json:"images"`
	Artists []ArtistSummaryResponse `json:"artists"`
	Match   int                     `json:"match"`
}

func AlbumResponseFromSummary(summary AlbumSummary) (AlbumSummaryResponse, error) {
	var zero AlbumSummaryResponse

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

	return AlbumSummaryResponse{
		ID:      summary.id.String(),
		Type:    summary.kind,
		Title:   summary.title,
		Images:  summary.images,
		Artists: artists,
		Match:   *summary.match,
	}, nil
}

func AlbumResponseFromSlice(summaries []AlbumSummary) ([]AlbumSummaryResponse, error) {
	var responses []AlbumSummaryResponse
	for _, summary := range summaries {
		response, err := AlbumResponseFromSummary(summary)
		if err != nil {
			return nil, err
		}
		responses = append(responses, response)
	}
	return responses, nil
}

type AlbumSummary struct {
	id      domain.ID
	title   string
	kind    domain.InfoType
	images  []domain.Image
	artists []ArtistSummary
	match   *int
}

func AlbumSummaryFromDeezer(source deezer.Track) (AlbumSummary, error) {
	id, err := domain.NewDeezerID(source.Album.ID)
	if err != nil {
		return AlbumSummary{}, err
	}

	return AlbumSummary{
		id:     id,
		title:  source.Album.Title,
		kind:   domain.AlbumType,
		images: []domain.Image{{URL: source.Album.Cover}},
	}, nil
}

func (a *AlbumSummary) SetArtists(artists []ArtistSummary) {
	a.artists = artists
}

// Matchable impl
func (a *AlbumSummary) MatchByQuery(query Query) {
	match := utils.RankString(query.String(), fmt.Sprintf(
		"%s - %s",
		a.title, a.artists[0].name,
	))
	a.match = &match
}

// SearchItem impl
func (a AlbumSummary) ID() domain.ID { return a.id }

func (a AlbumSummary) Name() string { return a.title }

func (a AlbumSummary) Kind() domain.InfoType { return a.kind }

func (a AlbumSummary) Images() []domain.Image { return a.images }

func (a AlbumSummary) Match() *int { return a.match }

func (a AlbumSummary) Artists() []ArtistSummary { return a.artists }
