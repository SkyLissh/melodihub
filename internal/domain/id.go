package domain

import (
	"errors"
	"fmt"
)

var ErrInvalidId = errors.New("invalid id")
var ErrInvalidProvider = errors.New("invalid provider")

type ID struct {
	provider string
	id       int
}

func NewId(provider string, id int) (*ID, error) {
	if provider == "" {
		return nil, ErrInvalidProvider
	}
	if id < 0 {
		return nil, ErrInvalidId
	}

	return &ID{provider, id}, nil
}

func (id *ID) String() string {
	return fmt.Sprintf("%s:%d", id.provider, id.id)
}

func (id *ID) Provider() string {
	return id.provider
}

func (id *ID) RawId() int {
	return id.id
}

func (id *ID) Id() string {
	return fmt.Sprintf("%s:%d", id.provider, id.id)
}
