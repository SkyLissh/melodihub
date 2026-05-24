package search

import (
	"errors"
	"strings"
	"unicode"
)

var (
	ErrInvalidQuery = errors.New("search query cannot be empty")
	ErrInvalidChar  = errors.New("search query cannot contain invalid characters")
)

type Query string

func ParseQuery(q string) (Query, error) {
	query := strings.TrimSpace(q)

	if query == "" {
		return "", ErrInvalidQuery
	}

	for _, char := range query {
		if unicode.IsLetter(char) || unicode.IsDigit(char) || unicode.IsSpace(char) {
			continue
		}

		return "", ErrInvalidChar
	}

	return Query(query), nil
}

func (q Query) String() string {
	return string(q)
}
