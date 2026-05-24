package album

import (
	"context"

	"github.com/skylissh/melodihub/internal/providers/deezer"
)

type AlbumDetailClient interface {
	GetDetail(context context.Context, id int) (*deezer.AlbumDetail, error)
}
