package deezer

import (
	"context"
	"fmt"

	"github.com/skylissh/melodihub/internal/validator"
)

type SearchService struct {
	provider  *Provider
	validator *validator.Validator
}

func (s *SearchService) Find(ctx context.Context, query string, limit uint) ([]Track, error) {
	client := s.provider.client
	var result Response[Data[[]Track]]

	_, err := client.R().
		SetContext(ctx).
		SetQueryParam("q", query).
		SetQueryParam("limit", fmt.Sprintf("%d", limit)).
		SetResult(&result).
		Get("search")

	if err != nil {
		return nil, ServerError(err)
	}

	if result.Error != nil {
		return nil, result.Error.ToError()
	}

	if err := s.validator.Validate(result.Data); err != nil {
		return nil, InvalidResponse(err)
	}

	for _, track := range result.Data.Data {
		if err := s.validator.Validate(track); err != nil {
			return nil, InvalidResponse(err)
		}
	}

	return result.Data.Data, nil
}
