package search

import "errors"

var ErrInvalidQuery = errors.New("search query cannot be empty")

type Query string

func ParseQuery(q string) (Query, error) {
	if q == "" {
		return "", ErrInvalidQuery
	}
	return Query(q), nil
}

func (q Query) String() string {
	return string(q)
}
