package search

type SeenData struct {
	Tracks  map[string]struct{}
	Albums  map[string]struct{}
	Artists map[string]struct{}
}

func NewSeenData() *SeenData {
	return &SeenData{
		Tracks:  make(map[string]struct{}),
		Albums:  make(map[string]struct{}),
		Artists: make(map[string]struct{}),
	}
}

func (sd *SeenData) MarkAlbum(albumID string) bool {
	if _, ok := sd.Albums[albumID]; !ok {
		sd.Albums[albumID] = struct{}{}
		return true
	}

	return false
}

func (sd *SeenData) MarkTrack(trackID string) bool {
	if _, ok := sd.Tracks[trackID]; !ok {
		sd.Tracks[trackID] = struct{}{}
		return true
	}

	return false
}

func (sd *SeenData) MarkArtist(artistID string) bool {
	if _, ok := sd.Artists[artistID]; !ok {
		sd.Artists[artistID] = struct{}{}
		return true
	}

	return false
}
