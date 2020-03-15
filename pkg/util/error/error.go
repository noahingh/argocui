package error

import "fmt"

var (
	// ErrNotImplement is the error when the method is not implemented.
	ErrNotImplement = fmt.Errorf("it's not implemented")
)
