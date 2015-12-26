package argjoy

import "errors"

var (
	NoMatchErr  = errors.New("Argument type has no matching codec")
	ArgCountErr = errors.New("Argument count mismatch")
)
