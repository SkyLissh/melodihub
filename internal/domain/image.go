// Package domain defines all common structs used by MelodiHub features.
package domain

type Image struct {
	URL    string `json:"url" validate:"required,http_url"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}
