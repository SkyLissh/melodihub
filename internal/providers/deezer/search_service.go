package deezer

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type SearchOrder string

const (
	SearchOrderRanking      SearchOrder = "RANKING"
	SearchOrderTrackAsc     SearchOrder = "TRACK_ASC"
	SearchOrderTrackDesc    SearchOrder = "TRACK_DESC"
	SearchOrderArtistAsc    SearchOrder = "ARTIST_ASC"
	SearchOrderArtistDesc   SearchOrder = "ARTIST_DESC"
	SearchOrderDurationAsc  SearchOrder = "DURATION_ASC"
	SearchOrderDurationDesc SearchOrder = "DURATION_DESC"
)

type SearchService interface {
	// TODO: Add more search parameters (e.g., limit, offset, order)
	Find(query string, limit int) ([]SearchResult, error)
}

type searchService struct {
	provider *Provider
}

func (s *searchService) Find(query string, limit int) ([]SearchResult, error) {
	client := s.provider.client
	var result Response[[]SearchResult]

	if limit <= 0 {
		return nil, errors.New("limit must be greater than 0")
	}

	_, err := client.R().
		SetQueryParam("q", query).
		SetQueryParam("limit", fmt.Sprintf("%d", limit)).
		SetResult(&result).
		Get("search")

	if err != nil {
		return nil, err
	}

	validate := validator.New()
	if err := validate.Struct(result); err != nil {
		return nil, err
	}

	return result.Data, nil
}
