package errors

import (
	"fmt"
	"github.com/username/rustx/token"
)

type SyntaxError struct {
	Message  string
	Position token.Position
}

func (e *SyntaxError) Error() string {
	if e.Position.Line > 0 {
		return fmt.Sprintf("SyntaxError: %s\nat line %d column %d", e.Message, e.Position.Line, e.Position.Column)
	}
	return fmt.Sprintf("SyntaxError: %s", e.Message)
}

func NewSyntaxError(pos token.Position, format string, args ...interface{}) *SyntaxError {
	return &SyntaxError{
		Message:  fmt.Sprintf(format, args...),
		Position: pos,
	}
}

type RuntimeError struct {
	Message  string
	Position token.Position
}

func (e *RuntimeError) Error() string {
	if e.Position.Line > 0 {
		return fmt.Sprintf("RuntimeError: %s\nat line %d column %d", e.Message, e.Position.Line, e.Position.Column)
	}
	return fmt.Sprintf("RuntimeError: %s", e.Message)
}

func NewRuntimeError(pos token.Position, format string, args ...interface{}) *RuntimeError {
	return &RuntimeError{
		Message:  fmt.Sprintf(format, args...),
		Position: pos,
	}
}
