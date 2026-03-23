package lastfm

type Tag struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Count int    `json:"count,omitempty"`
}

type TopTags struct {
	Tags []Tag `json:"tag"`
}
