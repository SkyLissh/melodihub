package spotify

type ArtistSimple struct {
	ID   string   `json:"id"`
	Name string   `json:"name"`
	URL  string   `json:"url"`
	Href string   `json:"href"`
	Type InfoType `json:"type"`
}

type Artist struct {
	ArtistSimple

	Popularity int       `json:"popularity"`
	Genres     []string  `json:"genres"`
	Images     []Image   `json:"images"`
	Followers  Followers `json:"followers"`
}
