package carousel

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"

	"github.com/skylissh/melodihub/internal/contracts"
	common "github.com/skylissh/melodihub/internal/domain"
	"github.com/skylissh/melodihub/internal/providers/deezer"
	"github.com/skylissh/melodihub/internal/providers/lastfm"
)

type lastfmGeoData struct {
	topArtists []lastfm.Artist
	topTracks  []lastfm.Track
}

type deezerGeoData struct {
	topArtists []deezer.SearchResult
	topTracks  []deezer.SearchResult
}

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
	country string,
) ([]Section, error) {

	if s.cache != nil {
		if cached, ok := s.cache.Get(ctx, country); ok {
			return cached, nil
		}
	}

	geoData, err := s.getLastfmGeoData(ctx, country)
	if err != nil {
		return nil, err
	}

	deezerData, err := s.getDeezerData(ctx, geoData)
	if err != nil {
		return nil, err
	}

	artists := make([]Item, len(deezerData.topArtists))
	for i, result := range deezerData.topArtists {
		artists[i] = Item{
			ID:   fmt.Sprintf("deezer:%d", result.Artist.ID),
			Type: common.ArtistType,
			Name: result.Artist.Name,
			Images: []common.Image{
				{
					URL: result.Artist.Picture,
				},
			},
		}
	}

	tracks := make([]Item, len(deezerData.topTracks))
	for i, result := range deezerData.topTracks {
		tracks[i] = Item{
			ID:   fmt.Sprintf("deezer:%d", result.ID),
			Type: common.AlbumType,
			Name: result.Title,
			Images: []common.Image{
				{
					URL: result.Album.Cover,
				},
			},
		}
	}

	carousels := []Section{}
	if len(artists) > 0 {
		carousels = append(carousels, Section{
			Title: "Top Artists",
			Items: artists,
		})
	}

	if len(tracks) > 0 {
		carousels = append(carousels, Section{
			Title: "Top Tracks",
			Items: tracks,
		})
	}

	if s.cache != nil {
		s.cache.Set(ctx, country, carousels)
	}

	return carousels, nil
}

func (s *Service) getLastfmGeoData(
	ctx context.Context,
	country string,
) (*lastfmGeoData, error) {
	g, gCtx := errgroup.WithContext(ctx)

	lastfmTopArtists := make([]lastfm.Artist, 0, 10)
	lastfmTopTracks := make([]lastfm.Track, 0, 10)

	g.Go(func() error {
		result, err := s.geo.GetTopArtists(gCtx, country, 10, 1)
		if err != nil {
			return err
		}

		lastfmTopArtists = result
		return nil
	})

	g.Go(func() error {
		result, err := s.geo.GetTopTracks(gCtx, country, 10, 1)
		if err != nil {
			return err
		}

		lastfmTopTracks = result
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return &lastfmGeoData{
		topArtists: lastfmTopArtists,
		topTracks:  lastfmTopTracks,
	}, nil
}

func (s *Service) getDeezerData(
	ctx context.Context,
	data *lastfmGeoData,
) (*deezerGeoData, error) {
	g, gCtx := errgroup.WithContext(ctx)
	g.SetLimit(5)

	deezerArtists := make([]deezer.SearchResult, len(data.topArtists))
	deezerTracks := make([]deezer.SearchResult, len(data.topTracks))

	for i, artist := range data.topArtists {
		g.Go(func() error {
			i, artist := i, artist // capture loop variables

			deezerArtist, err := s.search.Find(gCtx, artist.Name, 10)
			if err != nil {
				return err
			}
			deezerArtists[i] = deezerArtist[0]
			return nil
		})
	}

	for i, track := range data.topTracks {
		g.Go(func() error {
			i, track := i, track // capture loop variables

			deezerTrack, err := s.search.Find(gCtx, track.Name, 10)
			if err != nil {
				return err
			}

			deezerTracks[i] = deezerTrack[0]
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return &deezerGeoData{
		topArtists: deezerArtists,
		topTracks:  deezerTracks,
	}, nil
}
