package artist

import (
	"github.com/skylissh/melodihub/internal/domain"
	"github.com/skylissh/melodihub/internal/providers/deezer"
)

type TopAlbumResponse struct {
	ID     string          `json:"id" validate:"required"`
	Kind   domain.InfoType `json:"kind" validate:"required"`
	Title  string          `json:"title" validate:"required"`
	Images []domain.Image  `json:"images"`
}

func ResponseFromTopAlbum(album *TopAlbum) TopAlbumResponse {
	return TopAlbumResponse{
		ID:     album.ID.String(),
		Kind:   album.Kind,
		Title:  album.Title,
		Images: album.Images,
	}
}

func ResponseFromTopAlbums(albums []TopAlbum) []TopAlbumResponse {
	responses := make([]TopAlbumResponse, 0, len(albums))
	for _, album := range albums {
		responses = append(responses, ResponseFromTopAlbum(&album))
	}
	return responses
}

type TopAlbum struct {
	ID     domain.ID
	Kind   domain.InfoType
	Title  string
	Images []domain.Image
}

func TopAlbumFromDeezerAlbum(source *deezer.Album) (TopAlbum, error) {
	id, err := domain.NewDeezerID(source.ID)
	if err != nil {
		return TopAlbum{}, err
	}
	return TopAlbum{
		ID:     id,
		Kind:   domain.AlbumType,
		Title:  source.Title,
		Images: []domain.Image{{URL: source.Cover}},
	}, nil
}
