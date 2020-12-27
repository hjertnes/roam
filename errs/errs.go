package errs

import "github.com/rotisserie/eris"

var (
	Duplicate = eris.New("Duplicate")
	NotFound  = eris.New("not found")
)
