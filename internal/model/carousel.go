package model

type CarouselArtist struct {
	Name   string  `json:"name"`
	ID     string  `json:"id"`
	Images []Image `json:"images,omitempty"`
}

type CarouselItem struct {
	Type    InfoType         `json:"type"`
	ID      string           `json:"id"`
	Name    string           `json:"name"`
	Images  []Image          `json:"images,omitempty"`
	Artists []CarouselArtist `json:"artists,omitempty"`
}

type Carousel struct {
	Title       string         `json:"title"`
	Description *string        `json:"description,omitempty"`
	Items       []CarouselItem `json:"items"`
}
