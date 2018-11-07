package storage

import "errors"

var (
	// ErrResourceExists provides an error for when, in most cases, a record's
	// unique identifier already exists in the system.
	ErrResourceExists = errors.New("resource conflict")
)
