package httperror

// HTTPError carries an HTTP status code and response body from an upstream API
// so the handler can forward them to the client.
type HTTPError struct {
	StatusCode  int
	Body        []byte
	ContentType string
}

func (e *HTTPError) Error() string {
	if len(e.Body) > 0 {
		return string(e.Body)
	}
	return "upstream API returned error"
}
