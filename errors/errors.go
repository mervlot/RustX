package errors

import "fmt"

// SyntaxError represents a syntax error during lexing or parsing.
type SyntaxError struct {
	Line    int
	Column  int
	Message string
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("SyntaxError: %s at line %d column %d", e.Message, e.Line, e.Column)
}

// RuntimeError represents an error during script evaluation.
type RuntimeError struct {
	Line    int
	Column  int
	Message string
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("RuntimeError: %s at line %d column %d", e.Message, e.Line, e.Column)
}
