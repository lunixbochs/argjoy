package argjoy

import "fmt"

type NoMatchErr struct {
	desc string
}

func (e NoMatchErr) Error() string {
	return fmt.Sprintf("Argument conversion (%T) has no matching codec", e.desc)
}

var NoMatch = &NoMatchErr{}

type ArgCountErr struct {
	have, want int
}

func (e ArgCountErr) Error() string {
	return fmt.Sprintf("Argument count mismatch. Have %d, want %d", e.have, e.want)
}
