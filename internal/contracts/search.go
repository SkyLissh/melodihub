package contracts

import (
	"context"

	"github.com/skylissh/melodihub/internal/providers/deezer"
)

type DeezerSearchClient interface {
	Find(ctx context.Context, query string, limit int) ([]deezer.Track, error)
}
