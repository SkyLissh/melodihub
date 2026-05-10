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

func (k Kind) IsValid() bool {
	switch k {
	case AlbumKind, ArtistKind, TrackKind:
		return true
	default:
		return false
	}
}

func (k *Kind) UnmarshalJSON(data []byte) error {
	var res string
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	kind := Kind(res)
	if !kind.IsValid() {
		return fmt.Errorf("invalid kind: %s", res)
	}

	*k = kind
	return nil
}
