package lastfm

type responseError struct {
	Error   int    `json:"error"`
	Message string `json:"message"`
}

var errorCodeToStatus = map[int]int{
	1:  400, // Invalid service
	2:  400, // Invalid Method
	3:  401, // Authentication Failed
	4:  400, // Invalid format
	5:  401, // Bad authentication token
	6:  400, // Invalid parameters
	7:  404, // Invalid resource specified
	8:  500, // Operation failed
	9:  401, // Invalid session key
	10: 400, // Invalid method signature supplied
	11: 503, // Temporary server issues
	13: 400, // Invalid username
	14: 400, // Invalid timestamp
	15: 403, // Deleted API Key
	16: 400, // Not enough content
	17: 503, // Service temporarily unavailable
	18: 401, // Login: User required
	19: 400, // Subsonic server required
	20: 410, // Legacy API method no longer active
	21: 403, // Scrobbles to bad/suspended account
	22: 400, // This error does not exist
	23: 403, // Open registration disabled
	24: 429, // Too many requests in short time
	25: 403, // API key permissions not granted
	26: 403, // Suspended API key
	27: 429, // Rate limited
	29: 503, // Offline
}

type ResponseError struct {
	StatusCode int
	Message    string
}

func (e *ResponseError) Error() string {
	return e.Message
}

func newResponseError(errorCode int, message string) *ResponseError {
	statusCode, ok := errorCodeToStatus[errorCode]
	if !ok {
		statusCode = 500 // Default to Internal Server Error if code is unknown
	}

	return &ResponseError{
		StatusCode: statusCode,
		Message:    message,
	}
}
