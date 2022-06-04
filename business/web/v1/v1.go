package v1

import (
	"errors"
	"fmt"
	"net/http"
)

func NewRequestError(err error, status int) error {
	return &RequestError{Err: err, Status: status}
}

// Request error

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

// Fields error

func NewFieldError(field, message string) FieldError {
	return FieldError{Field: field, Message: message}
}

type FieldError struct {
	Field   string
	Message string
}

func (fe FieldError) Status() int {
	return http.StatusUnprocessableEntity
}

func (fe FieldError) Error() string {
	return fmt.Sprintf("%s %s", fe.Field, fe.Message)
}

func IsFieldError(err error) bool {
	var fe FieldError
	return errors.As(err, &fe)
}

func GetFieldError(err error) FieldError {
	var fe FieldError
	if !errors.As(err, &fe) {
		return FieldError{}
	}
	return fe
}

// Response

type ErrorResponse struct {
	Err     string            `json:"error"`
	Status  int               `json:"status"`
	Details map[string]string `json:"details,omitempty"`
}
