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
	return fmt.Sprintf("SyntaxError: %s at line %d column %d", e.Message, e.Pos.Line, e.Pos.Column)
}

func NewSyntaxError(pos token.Position, msg string) *SyntaxError {
	return &SyntaxError{
		Message: msg,
		Pos:     pos,
	}
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

func NewRuntimeError(pos token.Position, msg string) *RuntimeError {
	return &RuntimeError{
		Message: msg,
		Pos:     pos,
	}
}
