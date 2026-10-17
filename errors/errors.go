package errors

import (
	"fmt"
)

// SyntaxError represents an error during lexical analysis or parsing.
type SyntaxError struct {
	Message string
	Line    int
	Column  int
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("SyntaxError: %s\nat line %d column %d", e.Message, e.Line, e.Column)
}

func NewSyntaxError(line, column int, format string, args ...interface{}) *SyntaxError {
	return &SyntaxError{
		Message: fmt.Sprintf(format, args...),
		Line:    line,
		Column:  column,
	}
}

// RuntimeError represents an error during script evaluation.
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

func NewRuntimeError(line, column int, format string, args ...interface{}) *RuntimeError {
	return &RuntimeError{
		Message: fmt.Sprintf(format, args...),
		Line:    line,
		Column:  column,
	}
}
