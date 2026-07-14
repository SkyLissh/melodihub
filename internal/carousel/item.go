package carousel

import (
	"github.com/skylissh/melodihub/internal/domain"
	"github.com/skylissh/melodihub/internal/providers/deezer"
)

type ArtistSummary struct {
	id   domain.ID
	kind domain.InfoType
	name string
}

type ArtistSummaryResponse struct {
	ID   string          `json:"id" validate:"required"`
	Kind domain.InfoType `json:"kind" validate:"required"`
	Name string          `json:"name" validate:"required"`
}

func ResponseFromArtistSummary(summary ArtistSummary) ArtistSummaryResponse {
	return ArtistSummaryResponse{
		ID:   summary.id.String(),
		Kind: summary.kind,
		Name: summary.name,
	}
}

func ResponseFromArtistSummaries(summaries []ArtistSummary) []ArtistSummaryResponse {
	responses := make([]ArtistSummaryResponse, 0, len(summaries))
	for _, summary := range summaries {
		responses = append(responses, ResponseFromArtistSummary(summary))
	}
	return responses
}

type Item struct {
	kind    domain.InfoType
	id      domain.ID
	name    string
	images  []domain.Image
	artists []ArtistSummary
}

func TopArtistFromDeezer(track deezer.Track) (Item, error) {
	id, err := domain.NewDeezerID(track.Artist.ID)
	if err != nil {
		return Item{}, err
	}
	return Item{
		kind:   domain.ArtistType,
		id:     id,
		name:   track.Artist.Name,
		images: []domain.Image{{URL: track.Artist.Picture}},
	}, nil
}

func TopTrackFromDeezer(track deezer.Track) (Item, error) {
	id, err := domain.NewDeezerID(track.ID)
	if err != nil {
		return Item{}, err
	}

	artistID, err := domain.NewDeezerID(track.Artist.ID)
	if err != nil {
		return Item{}, err
	}

	return Item{
		kind:   domain.TrackType,
		id:     id,
		name:   track.Title,
		images: []domain.Image{{URL: track.Album.Cover}},
		artists: []ArtistSummary{{
			id:   artistID,
			kind: domain.ArtistType,
			name: track.Artist.Name,
		}},
	}, nil
}

type ItemResponse struct {
	Kind    domain.InfoType         `json:"kind"`
	ID      string                  `json:"id"`
	Name    string                  `json:"name"`
	Images  []domain.Image          `json:"images,omitempty"`
	Artists []ArtistSummaryResponse `json:"artists,omitempty"`
}

func ResponseFromItem(item Item) ItemResponse {
	return ItemResponse{
		ID:      item.id.String(),
		Kind:    item.kind,
		Name:    item.name,
		Images:  item.images,
		Artists: ResponseFromArtistSummaries(item.artists),
	}
}

func ResponseFromItems(items []Item) []ItemResponse {
	responses := make([]ItemResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, ResponseFromItem(item))
	}
	return responses
}
