package deezer

import "github.com/go-playground/validator/v10"

type ArtistService interface {
	GetByID(id string) (*Artist, error)
}

type artistService struct {
	provider *Provider
}

func (s *artistService) GetByID(id string) (*Artist, error) {
	client := s.provider.client
	result := &Artist{}

	_, err := client.R().
		SetPathParam("id", id).
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
