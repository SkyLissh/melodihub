package carousel

type Section struct {
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	Items       []Item  `json:"items"`
}
