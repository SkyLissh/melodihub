package domain

import (
	"errors"
	"fmt"
)

var ErrInvalidId = errors.New("invalid id")

type ID struct {
	provider Provider
	id       int
}

func NewID(provider string, id int) (ID, error) {
	var zero ID

	if id < 0 {
		return zero, ErrInvalidId
	}

	p, err := ParseProvider(provider)
	if err != nil {
		return zero, err
	}

	return ID{p, id}, nil
}

func NewDeezerID(id int) (ID, error) {
	return NewID("deezer", id)
}

func NewLastFMID(id int) (ID, error) {
	return NewID("lastfm", id)
}

func (id ID) String() string {
	return fmt.Sprintf("%s:%d", id.provider, id.id)
}

func (id ID) Provider() Provider {
	return id.provider
}

func (id ID) Id() int {
	return id.id
}
