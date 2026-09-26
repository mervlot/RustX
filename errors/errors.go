package errors

import (
	"fmt"
)

// SyntaxError represents an error during lexing or parsing.
type SyntaxError struct {
	Message string
	Line    int
	Column  int
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("SyntaxError: %s at line %d column %d", e.Message, e.Line, e.Column)
}

// NewSyntaxError returns a new SyntaxError pointer.
func NewSyntaxError(msg string, line, col int) *SyntaxError {
	return &SyntaxError{
		Message: msg,
		Line:    line,
		Column:  col,
	}
}

// RuntimeError represents an error during evaluation/runtime.
type RuntimeError struct {
	Message string
	Line    int
	Column  int
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("RuntimeError: %s at line %d column %d", e.Message, e.Line, e.Column)
}

// NewRuntimeError returns a new RuntimeError pointer.
func NewRuntimeError(msg string, line, col int) *RuntimeError {
	return &RuntimeError{
		Message: msg,
		Line:    line,
		Column:  col,
	}
}
