// Package mapper defines all mappers used by the common Melodi API
package mapper

import (
	"fmt"

	"github.com/skylissh/melodihub/internal/model"
	"github.com/skylissh/melodihub/internal/providers/deezer"
	"github.com/skylissh/melodihub/internal/utils"
)

func SearchFromDeezer(
	deezerResult []deezer.SearchResult,
	query *string,
) *model.SearchResult {
	tracks := []model.TrackResult{}
	artists := []model.ArtistResult{}
	albums := []model.AlbumResult{}

	trackSetID := make(map[string]struct{}, len(deezerResult))
	artistSetID := make(map[string]struct{}, len(deezerResult))
	albumSetID := make(map[string]struct{}, len(deezerResult))

	for _, result := range deezerResult {
		var artistMatch int
		if query != nil {
			artistMatch = utils.RankString(*query, result.Artist.Name)
		}
		artist := model.ArtistResult{
			ID:    fmt.Sprintf("deezer:%d", result.Artist.ID),
			Name:  result.Artist.Name,
			Image: result.Artist.Picture,
			Match: artistMatch,
			Type:  model.ArtistType,
		}

		if _, exists := artistSetID[artist.ID]; !exists {
			artists = append(artists, artist)
			artistSetID[artist.ID] = struct{}{}
		}

		var trackMatch int
		if query != nil {
			trackMatch = utils.RankString(*query, fmt.Sprintf(
				"%s - %s",
				result.Title, result.Artist.Name,
			))
		}

		track := model.TrackResult{
			ID:             fmt.Sprintf("deezer:%d", result.ID),
			Title:          result.Title,
			Duration:       result.Duration,
			ExplicitLyrics: result.ExplicitLyrics,
			Preview:        result.Preview,
			Match:          trackMatch,
			Rank:           result.Rank,
			Artist:         artist,
			Type:           model.TrackType,
			Image:          result.Album.Cover,
		}
		if _, exists := trackSetID[track.ID]; !exists {
			tracks = append(tracks, track)
			trackSetID[track.ID] = struct{}{}
		}

		var albumMatch int
		if query != nil {
			albumMatch = utils.RankString(*query, fmt.Sprintf(
				"%s - %s",
				result.Album.Title, result.Artist.Name,
			))
		}

		album := model.AlbumResult{
			ID:     fmt.Sprintf("deezer:%d", result.Album.ID),
			Title:  result.Album.Title,
			Image:  result.Album.Cover,
			Match:  albumMatch,
			Artist: artist,
			Type:   model.AlbumType,
		}

		if _, exists := albumSetID[album.ID]; !exists {
			albums = append(albums, album)
			albumSetID[album.ID] = struct{}{}
		}
	}

	return &model.SearchResult{
		Tracks:  tracks,
		Artists: artists,
		Albums:  albums,
	}
}
