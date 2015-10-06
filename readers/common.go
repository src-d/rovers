package readers

import "errors"

var (
	ErrUnexpectedStatusCode = errors.New("unexpected status code")
	ErrRateLimitExceeded    = errors.New("rate limit exceeded")
	ErrPartialResponse      = errors.New("received partial data")
)
