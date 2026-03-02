package lastfm

type SimilarArtistAttr struct {
	Artist string `json:"artist"`
}

type SimilarArtistContainer struct {
	Artists []Artist          `json:"artist"`
	Attr    SimilarArtistAttr `json:"@attr"`
}

type SimilarArtists struct {
	SimilarArtists SimilarArtistContainer `json:"similarartists"`
}
