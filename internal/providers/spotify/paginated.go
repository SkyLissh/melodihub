package spotify

type Paginated[T any] struct {
	Items    []T     `json:"items"`
	Href     string  `json:"href"`
	Limit    int     `json:"limit"`
	Offset   int     `json:"offset"`
	Total    int     `json:"total"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
}
