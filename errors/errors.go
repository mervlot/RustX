package errors

import (
	"fmt"
	"rustx/token"
)

// SyntaxError represents a parsing or lexical error with source coordinates.
type SyntaxError struct {
	Message  string
	Position token.Position
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("SyntaxError:\n%s\nat line %d column %d", e.Message, e.Position.Line, e.Position.Column)
}

// NewSyntaxError creates a new SyntaxError.
func NewSyntaxError(pos token.Position, format string, a ...interface{}) *SyntaxError {
	return &SyntaxError{
		Message:  fmt.Sprintf(format, a...),
		Position: pos,
	}
}

// RuntimeError represents an error during program evaluation.
type RuntimeError struct {
	Message  string
	Position token.Position
}

func (e *RuntimeError) Error() string {
	if e.Position.Line > 0 {
		return fmt.Sprintf("RuntimeError:\n%s\nat line %d column %d", e.Message, e.Position.Line, e.Position.Column)
	}
	return fmt.Sprintf("RuntimeError:\n%s", e.Message)
}

// NewRuntimeError creates a new RuntimeError.
func NewRuntimeError(pos token.Position, format string, a ...interface{}) *RuntimeError {
	return &RuntimeError{
		Message:  fmt.Sprintf(format, a...),
		Position: pos,
	}
}
