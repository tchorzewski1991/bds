package v1

import "errors"

func NewRequestError(err error, status int) error {
	return &RequestError{Err: err, Status: status}
}

// Request

type RequestError struct {
	Err    error
	Status int
}

func (re *RequestError) Error() string {
	return re.Err.Error()
}

func IsRequestError(err error) bool {
	var re *RequestError
	return errors.As(err, &re)
}

// GetRequestError returns a copy of the RequestError pointer.
func GetRequestError(err error) *RequestError {
	var re *RequestError
	if !errors.As(err, &re) {
		return nil
	}
	return re
}

// Response

type ErrorResponse struct {
	Err     string            `json:"error"`
	Status  int               `json:"status"`
	Details map[string]string `json:"details,omitempty"`
}
