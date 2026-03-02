package model

import (
	"encoding/json"
	"fmt"
)

type InfoType string

const (
	InfoTypeAlbum  InfoType = "album"
	InfoTypeArtist InfoType = "artist"
	InfoTypeTrack  InfoType = "track"
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
	case InfoTypeAlbum, InfoTypeArtist, InfoTypeTrack:
		*i = InfoType(value)
		return nil
	default:
		return fmt.Errorf("invalid info type: %s", value)
	}
}
