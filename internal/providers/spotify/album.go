package spotify

import (
	"encoding/json"
	"fmt"
)

type DatePrecision string

const (
	DatePrecisionYear  DatePrecision = "year"
	DatePrecisionMonth DatePrecision = "month"
	DatePrecisionDay   DatePrecision = "day"
)

func (d *DatePrecision) String() string {
	return string(*d)
}

func (d *DatePrecision) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	switch DatePrecision(value) {
	case DatePrecisionYear, DatePrecisionMonth, DatePrecisionDay:
		*d = DatePrecision(value)
		return nil
	default:
		return fmt.Errorf("invalid release precision: %s", value)
	}
}

type AlbumType string

const (
	AlbumTypeAlbum       AlbumType = "album"
	AlbumTypeSingle      AlbumType = "single"
	AlbumTypeCompilation AlbumType = "compilation"
)

func (a *AlbumType) String() string {
	return string(*a)
}

func (a *AlbumType) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	switch AlbumType(value) {
	case AlbumTypeAlbum, AlbumTypeSingle, AlbumTypeCompilation:
		*a = AlbumType(value)
		return nil
	default:
		return fmt.Errorf("invalid album type: %s", value)
	}
}

type AlbumSimple struct {
	ID                   string         `json:"id"`
	Name                 string         `json:"name"`
	URL                  string         `json:"url"`
	Href                 string         `json:"href"`
	Images               []Image        `json:"images"`
	ReleaseDate          string         `json:"release_date"`
	ReleaseDatePrecision DatePrecision  `json:"release_date_precision"`
	TotalTracks          int            `json:"total_tracks"`
	Artists              []ArtistSimple `json:"artists"`
	AlbumType            AlbumType      `json:"album_type"`
	Type                 InfoType       `json:"type"`
}

type Album struct {
	AlbumSimple

	Tracks Paginated[TrackSimple] `json:"tracks"`
}
