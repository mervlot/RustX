package errors

import "fmt"

// SyntaxError represents an error during lexing or parsing.
type SyntaxError struct {
	Message string
	Line    int
	Column  int
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("SyntaxError: %s\nat line %d column %d", e.Message, e.Line, e.Column)
}

// RuntimeError represents an error during interpretation/runtime.
type RuntimeError struct {
	Message string
	Line    int
	Column  int
}

func (e *RuntimeError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf("RuntimeError: %s\nat line %d column %d", e.Message, e.Line, e.Column)
	}
	return fmt.Sprintf("RuntimeError: %s", e.Message)
}
