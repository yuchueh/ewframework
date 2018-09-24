package utils

import "errors"

var (
	// ErrAbort custom error when user stop request handler manually.
	ErrAbort = errors.New("User stop run")
)
