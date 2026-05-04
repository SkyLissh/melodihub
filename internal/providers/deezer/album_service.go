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
	result := &AlbumDetail{}

	_, err := client.R().
		SetContext(context).
		SetPathParam("id", fmt.Sprintf("%d", id)).
		SetResult(result).
		Get("/album/{id}")

	if err != nil {
		return nil, err
	}

	if err := s.validator.Validate(result); err != nil {
		return nil, InvalidResponse(err)
	}

	return result, nil
}
