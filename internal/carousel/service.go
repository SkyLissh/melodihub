package carousel

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/skylissh/melodihub/internal/contracts"
	"github.com/skylissh/melodihub/internal/domain"
	"github.com/skylissh/melodihub/internal/providers/lastfm"
)

type Service struct {
	geo    LastfmGeoClient
	search contracts.DeezerSearchClient
	cache  *Cache
}

func NewService(
	geo LastfmGeoClient,
	search contracts.DeezerSearchClient,
	cache *Cache,
) *Service {
	return &Service{geo, search, cache}
}

func (s *Service) GetCarousels(
	ctx context.Context,
	country Country,
) ([]CarouselResponse, error) {

	if s.cache != nil {
		if cached, ok := s.cache.Get(ctx, country); ok {
			return cached, nil
		}
	}

	topArtists, topTracks, err := s.getTopGeoData(ctx, country)
	if err != nil {
		return nil, err
	}

	if len(topArtists) == 0 && len(topTracks) == 0 {
		return nil, domain.NotFound(errors.New("no top artists or tracks found"))
	}

	artists, tracks, err := s.topItemsFromGeoData(ctx, topArtists, topTracks)
	if err != nil {
		return nil, err
	}

	carousels := []Carousel{}
	if len(artists) > 0 {
		carousels = append(carousels, NewCarousel("Top Artists", "", artists))
	}

	if len(tracks) > 0 {
		carousels = append(carousels, NewCarousel("Top Tracks", "", tracks))
	}

	response := ResponseFromCarousels(carousels)

	if s.cache != nil {
		s.cache.Set(ctx, country, response)
	}

	return response, nil
}

func (s *Service) getTopGeoData(
	ctx context.Context,
	country Country,
) ([]lastfm.Artist, []lastfm.Track, error) {
	g, gCtx := errgroup.WithContext(ctx)

	lastfmTopArtists := make([]lastfm.Artist, 0, 10)
	lastfmTopTracks := make([]lastfm.Track, 0, 10)

	g.Go(func() error {
		result, err := s.geo.GetTopArtists(gCtx, country.String(), 10, 1)
		if err != nil {
			return CarouselErrorFromLastfmError(err)
		}

		lastfmTopArtists = result
		return nil
	})

	g.Go(func() error {
		result, err := s.geo.GetTopTracks(gCtx, country.String(), 10, 1)
		if err != nil {
			return CarouselErrorFromLastfmError(err)
		}

		lastfmTopTracks = result
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, nil, err
	}

	return lastfmTopArtists, lastfmTopTracks, nil
}

func (s *Service) topItemsFromGeoData(
	ctx context.Context,
	artists []lastfm.Artist,
	tracks []lastfm.Track,
) ([]Item, []Item, error) {
	g, gCtx := errgroup.WithContext(ctx)
	g.SetLimit(5)

	var artistMu sync.Mutex
	var trackMu sync.Mutex
	topArtists := make([]Item, 0, len(artists))
	topTracks := make([]Item, 0, len(tracks))

	for _, artist := range artists {
		g.Go(func() error {
			deezerArtist, err := s.search.Find(gCtx, artist.Name, 10)
			if err != nil {
				return CarouselErrorFromDeezerError(err)
			}

			if len(deezerArtist) == 0 {
				return nil
			}

			item, err := TopArtistFromDeezer(deezerArtist[0])
			if err != nil {
				return nil
			}

			artistMu.Lock()
			topArtists = append(topArtists, item)
			artistMu.Unlock()
			return nil
		})
	}

	for _, track := range tracks {
		g.Go(func() error {
			deezerTrack, err := s.search.Find(gCtx, track.Name, 10)
			if err != nil {
				return CarouselErrorFromDeezerError(err)
			}

			if len(deezerTrack) == 0 {
				return nil
			}

			item, err := TopTrackFromDeezer(deezerTrack[0])
			if err != nil {
				return nil
			}

			trackMu.Lock()
			topTracks = append(topTracks, item)
			trackMu.Unlock()
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, nil, err
	}

	return topArtists, topTracks, nil
}
