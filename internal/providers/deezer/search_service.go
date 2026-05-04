package deezer

import (
	"context"
	"errors"
	"fmt"

	"github.com/skylissh/melodihub/internal/validator"
)

type SearchService struct {
	provider  *Provider
	validator *validator.Validator
}

func (s *SearchService) Find(ctx context.Context, query string, limit int) ([]Track, error) {
	client := s.provider.client
	var result Response[Data[[]Track]]

	if limit <= 0 {
		return nil, errors.New("limit must be greater than 0")
	}

	_, err := client.R().
		SetContext(ctx).
		SetQueryParam("q", query).
		SetQueryParam("limit", fmt.Sprintf("%d", limit)).
		SetResult(&result).
		Get("search")

	if err != nil {
		return nil, ServerError("failed to search tracks")
	}

	if err := s.validator.Validate(result); err != nil {
		return nil, InvalidResponse(err)
	}

	return result.Data.Data, nil
}
