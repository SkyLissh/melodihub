package artist

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"

	"github.com/skylissh/melodihub/internal/contracts"
	common "github.com/skylissh/melodihub/internal/domain"
	"github.com/skylissh/melodihub/internal/providers/lastfm"
	"github.com/skylissh/melodihub/internal/search"
)

type lastfmArtistData struct {
	Detail    *lastfm.ArtistDetail
	TopTracks []lastfm.Track
	TopAlbums []lastfm.Album
	Similar   []lastfm.Artist
}

type artistData struct {
	similar   []Summary
	topTracks []TopTrack
	topAlbums []TopAlbum
}

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
	id string,
) (*Detail, error) {
	if s.cache != nil {
		if cached, ok := s.cache.Get(ctx, id); ok {
			return cached, nil
		}
	}

	deezerId, err := strconv.Atoi(strings.TrimPrefix(id, "deezer:"))
	if err != nil {
		return nil, fmt.Errorf("Invalid artist ID format: %w", err)
	}

	artist, err := s.artist.GetByID(ctx, deezerId)
	if err != nil {
		return nil, err
	}

	lastfmData, err := s.getLastfmData(ctx, artist.Name)
	if err != nil {
		return nil, err
	}

	data, err := s.getArtistData(ctx, artist.Name, lastfmData)
	if err != nil {
		return nil, err
	}

	artists := &Detail{
		Summary: Summary{
			ID:     fmt.Sprintf("deezer:%d", artist.ID),
			Name:   artist.Name,
			Images: []common.Image{{URL: artist.Picture}},
			Type:   common.ArtistType,
		},
		Listeners: lastfmData.Detail.Stats.Listeners,
		Bio:       lastfmData.Detail.Bio.Summary,
		Similar:   data.similar,
		TopTracks: data.topTracks,
		TopAlbums: data.topAlbums,
	}

	if s.cache != nil {
		s.cache.Set(ctx, id, artists)
	}

	return artists, nil
}

func (s *Service) getLastfmData(ctx context.Context, name string) (*lastfmArtistData, error) {
	g, gCtx := errgroup.WithContext(ctx)
	g.SetLimit(5)

	lastfmArtist := &lastfm.ArtistDetail{}
	lastfmTopTracks := make([]lastfm.Track, 0, 10)
	lastfmTopAlbums := make([]lastfm.Album, 0, 15)
	lastfmSimilar := make([]lastfm.Artist, 0, 5)

	g.Go(func() error {
		detail, err := s.lastfm.GetDetail(gCtx, name)
		if err != nil {
			return err
		}

		*lastfmArtist = *detail
		return nil
	})

	g.Go(func() error {
		tracks, err := s.lastfm.GetTopTracks(gCtx, name, 10)
		if err != nil {
			return err
		}

		lastfmTopTracks = tracks.Tracks
		return nil
	})

	g.Go(func() error {
		albums, err := s.lastfm.GetTopAlbums(gCtx, name, 15)
		if err != nil {
			return err
		}

		lastfmTopAlbums = albums.Album
		return nil
	})

	g.Go(func() error {
		similar, err := s.lastfm.GetSimilar(gCtx, name, 5)
		if err != nil {
			return err
		}

		lastfmSimilar = similar
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("failed to fetch artist details from Last.fm: %w", err)
	}

	return &lastfmArtistData{
		Detail:    lastfmArtist,
		TopTracks: lastfmTopTracks,
		TopAlbums: lastfmTopAlbums,
		Similar:   lastfmSimilar,
	}, nil
}

func (s *Service) getArtistData(
	ctx context.Context,
	name string,
	data *lastfmArtistData,
) (*artistData, error) {
	g, gCtx := errgroup.WithContext(ctx)
	g.SetLimit(5)

	similar := make([]Summary, 0, 5)
	topTracks := make([]TopTrack, 0, 10)
	topAlbums := make([]TopAlbum, 0, 15)

	g.Go(func() error {
		result, err := s.search.Find(gCtx, name, 100)
		if err != nil {
			return err
		}

		searchResult := search.ResultFromDeezer(search.Query(name), result)

		for _, a := range searchResult.Albums() {
			if len(topAlbums) >= 15 {
				break
			}

			_, ok := lo.Find(data.TopAlbums, func(la lastfm.Album) bool {
				return strings.EqualFold(la.Name, a.Name())
			})
			if !ok {
				continue
			}

			topAlbums = append(topAlbums, TopAlbum{
				ID:     a.ID().String(),
				Type:   common.AlbumType,
				Title:  a.Name(),
				Images: a.Images(),
			})
		}

		for _, t := range searchResult.Tracks() {
			if len(topTracks) >= 10 {
				break
			}

			m, ok := lo.Find(data.TopTracks, func(lt lastfm.Track) bool {
				return strings.EqualFold(lt.Name, t.Name())
			})
			if !ok {
				continue
			}

			topTracks = append(topTracks, TopTrack{
				ID:        t.ID().String(),
				Type:      common.TrackType,
				Title:     t.Name(),
				Duration:  t.Duration(),
				Images:    t.Images(),
				Listeners: m.Listeners,
			})
		}

		for _, a := range searchResult.Artists() {
			if len(similar) >= 5 {
				break
			}

			_, ok := lo.Find(data.Similar, func(r lastfm.Artist) bool {
				return strings.EqualFold(r.Name, a.Name())
			})
			if !ok {
				continue
			}

			similar = append(similar, NewSummary(
				a.ID().String(),
				a.Name(),
				a.Images(),
			))
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("failed to fetch artist data from Deezer: %w", err)
	}

	return &artistData{
		similar:   similar,
		topTracks: topTracks,
		topAlbums: topAlbums,
	}, nil
}
