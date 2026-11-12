package errors

import "fmt"

// SyntaxError represents an error during lexing or parsing.
type SyntaxError struct {
	Message string
	Line    int
	Column  int
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("SyntaxError:\n%s\nat line %d column %d", e.Message, e.Line, e.Column)
}

// NewSyntaxError creates a new SyntaxError.
func NewSyntaxError(msg string, line, col int) *SyntaxError {
	return &SyntaxError{
		Message: msg,
		Line:    line,
		Column:  col,
	}
}

// RuntimeError represents an error during interpretation/evaluation.
type RuntimeError struct {
	Message string
	Line    int
	Column  int
}

func (e *RuntimeError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf("RuntimeError:\n%s\nat line %d column %d", e.Message, e.Line, e.Column)
	}
	return fmt.Sprintf("RuntimeError:\n%s", e.Message)
}

// NewRuntimeError creates a new RuntimeError.
func NewRuntimeError(msg string, line, col int) *RuntimeError {
	return &RuntimeError{
		Message: msg,
		Line:    line,
		Column:  col,
	}
}
