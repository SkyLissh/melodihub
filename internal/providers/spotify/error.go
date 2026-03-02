package spotify

type ResponseError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (e *ResponseError) Error() string {
	return e.Message
}
