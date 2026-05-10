package deezer

import (
	"context"
	"fmt"

	"github.com/skylissh/melodihub/internal/validator"
)

type ArtistService struct {
	provider  *Provider
	validator *validator.Validator
}

func (s *ArtistService) GetByID(ctx context.Context, id int) (*Artist, error) {
	client := s.provider.client
	var result Response[Artist]

	_, err := client.R().
		SetContext(ctx).
		SetPathParam("id", fmt.Sprintf("%d", id)).
		SetResult(&result).
		Get("artist/{id}")

	if err != nil {
		return nil, ServerError(err)
	}

	if result.Error != nil {
		return nil, result.Error.ToError()
	}

	if err := s.validator.Validate(result.Data); err != nil {
		return nil, InvalidResponse(err)
	}

	return result.Data, nil
}
