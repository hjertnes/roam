// Package errs contain errors
package errs

import "github.com/rotisserie/eris"

var (
	// ErrDuplicate is the error for a duplicate.
	ErrDuplicate = eris.New("duplicate")
	// ErrNotFound is the error for not found.
	ErrNotFound = eris.New("not found")
	// ErrNoValue is the error for when no value is found.
	ErrNoValue = eris.New("no value")
	ErrNever = eris.New("never ran")
	ErrTooMany = eris.New("too many matches")
	ErrNoop = eris.New("noop")
)
