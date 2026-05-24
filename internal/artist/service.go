package artist

import (
	"context"
	"strings"

	"golang.org/x/sync/errgroup"

	"github.com/skylissh/melodihub/internal/contracts"
	"github.com/skylissh/melodihub/internal/domain"
	"github.com/skylissh/melodihub/internal/providers/deezer"
	"github.com/skylissh/melodihub/internal/providers/lastfm"
)

type Service struct {
	artist GetArtistClient
	search contracts.DeezerSearchClient
	lastfm GetArtistDetailClient
	cache  *Cache
}

func NewService(
	artist GetArtistClient,
	search contracts.DeezerSearchClient,
	lastfm GetArtistDetailClient,
	cache *Cache,
) *Service {
	return &Service{
		artist,
		search,
		lastfm,
		cache,
	}
}

func (s *Service) GetArtistDetails(
	ctx context.Context,
	id domain.ID,
) (*DetailResponse, error) {
	if s.cache != nil {
		if cached, ok := s.cache.Get(ctx, id); ok {
			return cached, nil
		}
	}

	deezerArtist, err := s.artist.GetByID(ctx, id.Id())
	if err != nil {
		return nil, ArtistErrorFromDeezerError(err)
	}

	summary, err := SummaryFromDeezer(*deezerArtist)
	if err != nil {
		return nil, err
	}

	deezerTracks, err := s.search.Find(ctx, deezerArtist.Name, 100)
	if err != nil {
		return nil, ArtistErrorFromDeezerError(err)
	}

	g, gCtx := errgroup.WithContext(ctx)
	g.SetLimit(4)

	var detail *lastfm.ArtistDetail
	var similar []Summary
	var topTracks []TopTrack
	var topAlbums []TopAlbum

	g.Go(func() error {
		d, err := s.lastfm.GetDetail(gCtx, deezerArtist.Name)
		if err != nil {
			return ArtistErrorFromLastfmError(err)
		}
		detail = d
		return nil
	})

	g.Go(func() error {
		var err error
		similar, err = s.fetchAndBuildSimilar(gCtx, deezerArtist.Name, deezerTracks)
		return err
	})

	g.Go(func() error {
		var err error
		topAlbums, err = s.fetchAndBuildTopAlbums(gCtx, deezerArtist.Name, deezerTracks)
		return err
	})

	g.Go(func() error {
		var err error
		topTracks, err = s.fetchAndBuildTopTracks(gCtx, deezerArtist.Name, deezerTracks)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	artist := Detail{
		Summary:   summary,
		Listeners: detail.Stats.Listeners,
		Bio:       detail.Bio.Summary,
		Similar:   similar,
		TopTracks: topTracks,
		TopAlbums: topAlbums,
	}

	response := ResponseFromDetail(&artist)

	if s.cache != nil {
		s.cache.Set(ctx, id, &response)
	}

	return &response, nil
}

func (s *Service) fetchAndBuildSimilar(
	ctx context.Context,
	name string,
	deezerTracks []deezer.Track,
) ([]Summary, error) {
	similarRaw, err := s.lastfm.GetSimilar(ctx, name, 10)
	if err != nil {
		return nil, ArtistErrorFromLastfmError(err)
	}
	return s.buildSimilarArtists(deezerTracks, similarRaw)
}

func (s *Service) fetchAndBuildTopAlbums(
	ctx context.Context,
	name string,
	deezerTracks []deezer.Track,
) ([]TopAlbum, error) {
	albums, err := s.lastfm.GetTopAlbums(ctx, name, 15)
	if err != nil {
		return nil, ArtistErrorFromLastfmError(err)
	}
	return s.buildTopAlbums(deezerTracks, albums.Album)
}

func (s *Service) fetchAndBuildTopTracks(
	ctx context.Context,
	name string,
	deezerTracks []deezer.Track,
) ([]TopTrack, error) {
	tracks, err := s.lastfm.GetTopTracks(ctx, name, 10)
	if err != nil {
		return nil, ArtistErrorFromLastfmError(err)
	}
	return s.buildTopTracks(deezerTracks, tracks.Tracks)
}

func (s *Service) buildSimilarArtists(
	tracks []deezer.Track,
	similar []lastfm.Artist,
) ([]Summary, error) {
	similarNames := make(map[string]struct{}, len(similar))
	for _, a := range similar {
		similarNames[strings.ToLower(a.Name)] = struct{}{}
	}

	result := make([]Summary, 0, 5)
	seen := make(map[string]struct{})

	for _, track := range tracks {
		if len(result) >= 5 {
			break
		}

		artistName := strings.ToLower(track.Artist.Name)
		if _, ok := similarNames[artistName]; !ok {
			continue
		}
		if _, ok := seen[artistName]; ok {
			continue
		}
		seen[artistName] = struct{}{}

		summary, err := SummaryFromDeezer(track.Artist)
		if err != nil {
			continue
		}

		result = append(result, summary)
	}

	return result, nil
}

func (s *Service) buildTopAlbums(
	tracks []deezer.Track,
	lastfmAlbums []lastfm.Album,
) ([]TopAlbum, error) {
	topAlbums := make([]TopAlbum, 0, 15)
	seen := make(map[string]struct{})

	for _, track := range tracks {
		if len(topAlbums) >= 15 {
			break
		}

		albumName := strings.ToLower(track.Album.Title)

		matched := false
		for _, la := range lastfmAlbums {
			if strings.EqualFold(la.Name, track.Album.Title) {
				matched = true
				break
			}
		}
		if !matched {
			continue
		}

		if _, ok := seen[albumName]; ok {
			continue
		}
		seen[albumName] = struct{}{}

		album, err := TopAlbumFromDeezerAlbum(&track.Album)
		if err != nil {
			continue
		}

		topAlbums = append(topAlbums, album)
	}

	return topAlbums, nil
}

func (s *Service) buildTopTracks(
	tracks []deezer.Track,
	lastfmTracks []lastfm.Track,
) ([]TopTrack, error) {
	topTracks := make([]TopTrack, 0, 10)
	seen := make(map[string]struct{})

	for _, track := range tracks {
		if len(topTracks) >= 10 {
			break
		}

		trackName := strings.ToLower(track.Title)

		var matchedTrack *lastfm.Track
		for i := range lastfmTracks {
			if strings.EqualFold(lastfmTracks[i].Name, track.Title) {
				matchedTrack = &lastfmTracks[i]
				break
			}
		}
		if matchedTrack == nil {
			continue
		}

		if _, ok := seen[trackName]; ok {
			continue
		}
		seen[trackName] = struct{}{}

		topTrack, err := TopTrackFromDeezer(&track)
		if err != nil {
			continue
		}

		topTrack.SetListeners(matchedTrack.Listeners)
		topTracks = append(topTracks, topTrack)
	}

	return topTracks, nil
}
