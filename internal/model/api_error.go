package model

type APIError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}
