package domain

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidId = errors.New("invalid id")
var ErrEmptyId = errors.New("empty id")

type ID struct {
	provider Provider
	id       int
}

func ParseID(s string) (ID, error) {
	var zero ID
	if s == "" {
		return zero, ErrEmptyId
	}

	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return zero, ErrInvalidId
	}

	provider := parts[0]
	id, err := strconv.Atoi(parts[1])
	if err != nil {
		return zero, ErrInvalidId
	}

	return NewID(provider, id)
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
