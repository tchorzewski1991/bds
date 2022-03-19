package web

import "errors"

func NewShutdownError(message string) error {
	return &shutdownError{message: message}
}

type shutdownError struct {
	message string
}

func (se *shutdownError) Error() string {
	return se.message
}

func IsShutdownError(err error) bool {
	var se *shutdownError
	return errors.As(err, &se)
}
