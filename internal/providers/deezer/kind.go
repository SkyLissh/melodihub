package deezer

import (
	"encoding/json"
	"fmt"
)

type Kind string

const (
	AlbumKind  Kind = "album"
	ArtistKind Kind = "artist"
	TrackKind  Kind = "track"
)

func (k Kind) String() string {
	return string(k)
}

func ParseKind(s string) (Kind, error) {
	switch s {
	case "album":
		return AlbumKind, nil
	case "artist":
		return ArtistKind, nil
	case "track":
		return TrackKind, nil
	default:
		return "", fmt.Errorf("unknown kind: %s", s)
	}
}

func (k Kind) IsValid() bool {
	switch k {
	case AlbumKind, ArtistKind, TrackKind:
		return true
	default:
		return false
	}
}

func (k Kind) MarshalJSON() ([]byte, error) {
	return []byte(`"` + k.String() + `"`), nil
}

func (k *Kind) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	parsed, err := ParseKind(s)
	if err != nil {
		return err
	}
	*k = parsed
	return nil
}
