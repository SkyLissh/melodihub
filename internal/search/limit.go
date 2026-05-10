package search

import (
	"errors"
	"strconv"
)

var (
	ErrLimitInvalid  = errors.New("search limit must be greater than 0")
	ErrLimitTooLarge = errors.New("search limit is too large (max 100)")
)

const DefaultLimit Limit = 10

type Limit uint

func ParseLimit(value string) (Limit, error) {
	if value == "" {
		return DefaultLimit, nil
	}

	val, err := strconv.ParseUint(value, 10, 0)
	if err != nil {
		return 0, ErrLimitInvalid
	}

	if val == 0 {
		return 0, ErrLimitInvalid
	}

	if val > 100 {
		return 0, ErrLimitTooLarge
	}

	return Limit(val), nil
}

func (l Limit) Value() uint {
	return uint(l)
}

func (l Limit) String() string {
	return strconv.FormatUint(uint64(l), 10)
}
