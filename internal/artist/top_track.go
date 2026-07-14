package artist

import (
	"github.com/skylissh/melodihub/internal/domain"
	"github.com/skylissh/melodihub/internal/providers/deezer"
)

type TopTrackResponse struct {
	ID        string          `json:"id" validate:"required"`
	Kind      domain.InfoType `json:"kind" validate:"required"`
	Title     string          `json:"title" validate:"required"`
	Images    []domain.Image  `json:"images"`
	Duration  int             `json:"duration" validate:"min=1"`
	Listeners int             `json:"listeners" validate:"min=1"`
}

func ResponseFromTopTrack(t *TopTrack) TopTrackResponse {
	return TopTrackResponse{
		ID:        t.id.String(),
		Kind:      t.kind,
		Title:     t.title,
		Images:    t.images,
		Duration:  t.duration,
		Listeners: t.listeners,
	}
}

func ResponseFromTopTracks(tracks []TopTrack) []TopTrackResponse {
	responses := make([]TopTrackResponse, 0, len(tracks))
	for _, track := range tracks {
		responses = append(responses, ResponseFromTopTrack(&track))
	}
	return responses
}

type TopTrack struct {
	id        domain.ID
	kind      domain.InfoType
	title     string
	images    []domain.Image
	duration  int
	listeners int
}

func TopTrackFromDeezer(t *deezer.Track) (TopTrack, error) {
	id, err := domain.NewDeezerID(t.ID)
	if err != nil {
		return TopTrack{}, err
	}

	return TopTrack{
		id:       id,
		kind:     domain.TrackType,
		title:    t.Title,
		images:   []domain.Image{{URL: t.Album.Cover}},
		duration: t.Duration,
	}, nil
}

func (t *TopTrack) SetListeners(listeners int) {
	t.listeners = listeners
}
