package errors

import "fmt"

// SyntaxError represents an error encountered during lexical analysis or parsing.
type SyntaxError struct {
	Message string
	Line    int
	Column  int
}

func (e *SyntaxError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf("SyntaxError:\n%s\nat line %d column %d", e.Message, e.Line, e.Column)
	}
	return fmt.Sprintf("SyntaxError:\n%s", e.Message)
}

// RuntimeError represents an error encountered during code execution.
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
