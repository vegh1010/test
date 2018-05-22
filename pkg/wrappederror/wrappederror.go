package wrappederror

import (
	"bytes"
	"fmt"
)

// Err is an implementation of the error interface, but contains
// an error and a description for it.
type Err struct {
	err         error
	description string
}

// NewErr creates a new *Err from an error and a description.
func NewErr(err error, description string) error {
	return &Err{err, description}
}

// Error implements the error interface for Err.
func (e *Err) Error() string {
	if e.err == nil {
		return fmt.Sprintf("<nil>: %s", e.description)
	}
	return fmt.Sprintf("%s: %s", e.description, e.err.Error())
}

// Errors is a slice of *Err.
type Errors []*Err

// Add adds an Err to the Errors slice.
func (e *Errors) Add(err error, description string) {
	*e = append(*e, &Err{err, description})
}

// AddErr adds an Err type to the Errors slice.
//
// If the conversion from the error type to an Err
// type occurs, AddErr will add an Err with the original
// error and a description that it was unable to be converted.
func (e *Errors) AddErr(err error) {
	errType, ok := err.(*Err)
	if !ok {
		e.Add(err, "Error converting error to *Err type")
		return
	}
	e.Add(errType.err, errType.description)
}

// Error returns a concatenated version of all errors contained within an Errors type.
func (e Errors) Error() string {
	var buf bytes.Buffer

	buf.WriteString("Errors: ")
	for i, err := range e {
		buf.WriteString(fmt.Sprintf("[%s]", err.Error()))
		if i < len(e)-1 {
			buf.WriteString("; ")
		}
	}

	return buf.String()
}

// New returns a new Errors variable.
//
// NOTE: If wanting to know if Errors is nil, it should be checked
//       with `if len(werrors) > 0` as opposed to `if werrors == nil`.
func New() Errors {
	var e Errors
	return e
}
