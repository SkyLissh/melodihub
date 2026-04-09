package search

import (
	"fmt"

	"github.com/skylissh/melodihub/internal/providers/deezer"
	"github.com/skylissh/melodihub/internal/utils"

	common "github.com/skylissh/melodihub/internal/domain"
)

func Builder(
	deezerResult []deezer.SearchResult,
	query *string,
) *Result {
	tracks := []TrackSummary{}
	artists := []ArtistSummary{}
	albums := []AlbumSummary{}

	trackSetID := make(map[string]struct{}, len(deezerResult))
	artistSetID := make(map[string]struct{}, len(deezerResult))
	albumSetID := make(map[string]struct{}, len(deezerResult))

	for _, result := range deezerResult {
		var artistMatch int
		if query != nil {
			artistMatch = utils.RankString(*query, result.Artist.Name)
		}
		artistData := ArtistSummary{
			ID:   fmt.Sprintf("deezer:%d", result.Artist.ID),
			Name: result.Artist.Name,
			Images: []common.Image{
				{
					URL: result.Artist.Picture,
				},
			},
			Match: artistMatch,
			Type:  common.ArtistType,
		}

		if _, exists := artistSetID[artistData.ID]; !exists {
			artists = append(artists, artistData)
			artistSetID[artistData.ID] = struct{}{}
		}

		var trackMatch int
		if query != nil {
			trackMatch = utils.RankString(*query, fmt.Sprintf(
				"%s - %s",
				result.Title, result.Artist.Name,
			))
		}

		track := TrackSummary{
			ID:       fmt.Sprintf("deezer:%d", result.ID),
			Title:    result.Title,
			Duration: result.Duration,
			Match:    trackMatch,
			Artists:  []ArtistSummary{artistData},
			Type:     common.TrackType,
			Images:   []common.Image{{URL: result.Album.Cover}},
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

		album := AlbumSummary{
			ID:      fmt.Sprintf("deezer:%d", result.Album.ID),
			Title:   result.Album.Title,
			Images:  []common.Image{{URL: result.Album.Cover}},
			Match:   albumMatch,
			Artists: []ArtistSummary{artistData},
			Type:    common.AlbumType,
		}

		if _, exists := albumSetID[album.ID]; !exists {
			albums = append(albums, album)
			albumSetID[album.ID] = struct{}{}
		}
	}

	results := &Result{
		Tracks:  tracks,
		Artists: artists,
		Albums:  albums,
	}

	results.FindTopResult()

	return results
}
