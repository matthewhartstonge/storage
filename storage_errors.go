package storage

import "errors"

var (
	// Provides an error for conflicting records.
	ErrResourceExists = errors.New("resource conflict")
)
