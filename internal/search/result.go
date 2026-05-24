package search

import (
	"fmt"
	"slices"

	"github.com/skylissh/melodihub/internal/domain"
	"github.com/skylissh/melodihub/internal/providers/deezer"
)

type SearchItem interface {
	ID() domain.ID
	Name() string
	Kind() domain.InfoType
	Images() []domain.Image
	Match() *int
	Artists() []ArtistSummary
}

type ResultResponse struct {
	TopResult TopResultResponse       `json:"top_result"`
	Tracks    []TrackSummaryResponse  `json:"tracks"`
	Artists   []ArtistSummaryResponse `json:"artists"`
	Albums    []AlbumSummaryResponse  `json:"albums"`
}

func ResponseFromResult(result *Result) (*ResultResponse, error) {
	if result == nil {
		return nil, fmt.Errorf("result cannot be nil")
	}

	tracks, err := TrackResponseFromSlice(result.tracks)
	if err != nil {
		return nil, err
	}

	artists, err := ArtistResponseFromSlice(result.artists)
	if err != nil {
		return nil, err
	}

	albums, err := AlbumResponseFromSlice(result.albums)
	if err != nil {
		return nil, err
	}

	topResult, err := result.FindTopResult()
	if err != nil {
		return nil, err
	}

	topResultRes, err := TopResultResponseFromSummary(topResult)
	if err != nil {
		return nil, err
	}

	return &ResultResponse{
		TopResult: topResultRes,
		Tracks:    tracks,
		Artists:   artists,
		Albums:    albums,
	}, nil
}

type Result struct {
	tracks  []TrackSummary
	artists []ArtistSummary
	albums  []AlbumSummary
}

func (r Result) Tracks() []TrackSummary {
	return r.tracks
}

func (r Result) Artists() []ArtistSummary {
	return r.artists
}

func (r Result) Albums() []AlbumSummary {
	return r.albums
}

func ResultFromDeezer(query Query, source []deezer.Track) *Result {
	result := &Result{
		tracks:  make([]TrackSummary, 0, len(source)),
		artists: make([]ArtistSummary, 0, len(source)),
		albums:  make([]AlbumSummary, 0, len(source)),
	}

	seen := NewSeenData()

	for _, src := range source {
		artistData, err := ArtistSummaryFromDeezer(src.Artist)
		if err != nil {
			continue
		}

		artistData.MatchByQuery(query)

		if seen.MarkArtist(artistData.ID().String()) {
			result.artists = append(result.artists, artistData)
		}

		track, err := TrackSummaryFromDeezer(src)
		if err != nil {
			continue
		}

		track.SetArtists([]ArtistSummary{artistData})
		track.MatchByQuery(query)

		if seen.MarkTrack(track.ID().String()) {
			result.tracks = append(result.tracks, track)
		}

		album, err := AlbumSummaryFromDeezer(src)
		if err != nil {
			continue
		}

		album.SetArtists([]ArtistSummary{artistData})
		album.MatchByQuery(query)

		if seen.MarkAlbum(album.ID().String()) {
			result.albums = append(result.albums, album)
		}
	}

	return result
}

func (r *Result) FindTopResult() (TopResult, error) {
	var zero TopResult

	length := len(r.tracks) + len(r.artists) + len(r.albums)
	if length == 0 {
		return zero, fmt.Errorf("no results found")
	}

	results := make([]SearchItem, 0, length)

	for _, track := range r.tracks {
		results = append(results, track)
	}
	for _, artist := range r.artists {
		results = append(results, artist)
	}
	for _, album := range r.albums {
		results = append(results, album)
	}

	highest := slices.MaxFunc(results, func(x, y SearchItem) int {
		return *x.Match() - *y.Match()
	})

	result, err := TopResultFromSearchItem(highest)
	if err != nil {
		return zero, err
	}

	return result, nil
}
