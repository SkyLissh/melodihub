package carousel

type Carousel struct {
	title       string
	description string
	items       []Item
}

func NewCarousel(title string, description string, items []Item) Carousel {
	return Carousel{
		title:       title,
		description: description,
		items:       items,
	}
}

type CarouselResponse struct {
	Title       string         `json:"title"`
	Description string         `json:"description,omitempty"`
	Items       []ItemResponse `json:"items"`
}

func ResponseFromCarousel(Carousel Carousel) CarouselResponse {
	return CarouselResponse{
		Title:       Carousel.title,
		Description: Carousel.description,
		Items:       ResponseFromItems(Carousel.items),
	}
}

func ResponseFromCarousels(Carousels []Carousel) []CarouselResponse {
	responses := make([]CarouselResponse, 0, len(Carousels))
	for _, Carousel := range Carousels {
		responses = append(responses, ResponseFromCarousel(Carousel))
	}
	return responses
}
