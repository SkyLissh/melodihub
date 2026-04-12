package domain

import (
	"encoding/json"
	"fmt"
)

type InfoType string

const (
	AlbumType  InfoType = "album"
	ArtistType InfoType = "artist"
	TrackType  InfoType = "track"
)

func (i *InfoType) String() string {
	return string(*i)
}

func (i *InfoType) UnmarshalJSON(data []byte) error {
	var value string

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	switch InfoType(value) {
	case AlbumType, ArtistType, TrackType:
		*i = InfoType(value)
		return nil
	default:
		return fmt.Errorf("invalid info type: %s", value)
	}
}
