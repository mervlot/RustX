package errors

import "fmt"

// SyntaxError represents a syntax or lexical error with line and column information.
type SyntaxError struct {
	Message string
	Line    int
	Column  int
	File    string
}

func (e *SyntaxError) Error() string {
	if e.File != "" {
		return fmt.Sprintf("SyntaxError in %s:\n%s\nat line %d column %d", e.File, e.Message, e.Line, e.Column)
	}
	return fmt.Sprintf("SyntaxError:\n%s\nat line %d column %d", e.Message, e.Line, e.Column)
}

// RuntimeError represents an error during script evaluation.
type RuntimeError struct {
	Message string
	Line    int
	Column  int
	File    string
}

func (e *RuntimeError) Error() string {
	if e.File != "" {
		return fmt.Sprintf("RuntimeError in %s:\n%s\nat line %d column %d", e.File, e.Message, e.Line, e.Column)
	}
	return fmt.Sprintf("RuntimeError:\n%s\nat line %d column %d", e.Message, e.Line, e.Column)
}
