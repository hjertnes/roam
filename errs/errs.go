package errs

import "github.com/rotisserie/eris"

var Duplicate = eris.New("Duplicate")
var NotFound = eris.New("not found")