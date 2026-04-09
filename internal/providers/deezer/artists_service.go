package deezer

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ArtistService interface {
	GetByID(ctx context.Context, id int) (*Artist, error)
}

type artistService struct {
	provider *Provider
}

func (s *artistService) GetByID(ctx context.Context, id int) (*Artist, error) {
	client := s.provider.client
	result := &Artist{}

	_, err := client.R().
		SetContext(ctx).
		SetPathParam("id", fmt.Sprintf("%d", id)).
		SetResult(result).
		Get("artist/{id}")

	if err != nil {
		return nil, err
	}

	validate := validator.New()
	if err := validate.Struct(result); err != nil {
		return nil, err
	}

	return result, nil
}
