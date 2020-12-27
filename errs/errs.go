// Package errs contain errors
package errs

import "github.com/rotisserie/eris"

var (
	// ErrDuplicate is the error for a duplicate.
	ErrDuplicate = eris.New("duplicate")
	// ErrNotFound is the error for not found.
	ErrNotFound = eris.New("not found")
)
