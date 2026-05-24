package artist

type DetailResponse struct {
	SummaryResponse

	Listeners int                `json:"listeners" validate:"required,min=10"`
	Bio       string             `json:"bio,omitempty"`
	TopTracks []TopTrackResponse `json:"top_tracks"`
	TopAlbums []TopAlbumResponse `json:"top_albums"`
	Similar   []SummaryResponse  `json:"similar_artists"`
}

func ResponseFromDetail(detail *Detail) DetailResponse {
	return DetailResponse{
		SummaryResponse: ResponseFromSummary(&detail.Summary),
		Listeners:       detail.Listeners,
		Bio:             detail.Bio,
		TopTracks:       ResponseFromTopTracks(detail.TopTracks),
		TopAlbums:       ResponseFromTopAlbums(detail.TopAlbums),
		Similar:         ResponseFromSummaries(detail.Similar),
	}
}

type Detail struct {
	Summary

	Listeners int
	Bio       string

	TopTracks []TopTrack
	TopAlbums []TopAlbum
	Similar   []Summary
}
