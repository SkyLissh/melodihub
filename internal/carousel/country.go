package carousel

import (
	"errors"
	"strings"
	"unicode"
)

var ErrCountryInvalid = errors.New("country must contain only letters and spaces")

const DefaultCountry Country = "united states"

type Country string

func ParseCountry(value string) (Country, error) {
	country := strings.ToLower(strings.TrimSpace(value))
	if country == "" {
		return DefaultCountry, nil
	}

	for _, char := range country {
		if unicode.IsLetter(char) || unicode.IsSpace(char) {
			continue
		}

		return "", ErrCountryInvalid
	}

	return Country(strings.Join(strings.Fields(country), " ")), nil
}

func (c Country) String() string {
	return string(c)
}
