package errors

import (
	"fmt"
	"github.com/username/rustx/token"
)

type SyntaxError struct {
	Message string
	Pos     token.Position
}

func (e *SyntaxError) Error() string {
	if e.Pos.Line > 0 {
		return fmt.Sprintf("SyntaxError: %s at line %d column %d", e.Message, e.Pos.Line, e.Pos.Column)
	}
	return fmt.Sprintf("SyntaxError: %s", e.Message)
}

type RuntimeError struct {
	Message string
	Pos     token.Position
}

func (e *RuntimeError) Error() string {
	if e.Pos.Line > 0 {
		return fmt.Sprintf("RuntimeError: %s at line %d column %d", e.Message, e.Pos.Line, e.Pos.Column)
	}
	return fmt.Sprintf("RuntimeError: %s", e.Message)
}

func NewSyntaxError(pos token.Position, message string) *SyntaxError {
	return &SyntaxError{
		Message: message,
		Pos:     pos,
	}
}

func NewRuntimeError(pos token.Position, message string) *RuntimeError {
	return &RuntimeError{
		Message: message,
		Pos:     pos,
	}
}
