package deezer

import (
	"context"
	"fmt"

	"github.com/skylissh/melodihub/internal/validator"
)

type AlbumService struct {
	provider  *Provider
	validator *validator.Validator
}

func (s *AlbumService) GetDetail(context context.Context, id int) (*AlbumDetail, error) {
	client := s.provider.client
	var result Response[AlbumDetail]

	_, err := client.R().
		SetContext(context).
		SetPathParam("id", fmt.Sprintf("%d", id)).
		SetResult(&result).
		Get("/album/{id}")

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
