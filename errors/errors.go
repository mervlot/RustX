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
		return fmt.Sprintf("SyntaxError:\n%s\nat line %d column %d", e.Message, e.Pos.Line, e.Pos.Column)
	}
	return fmt.Sprintf("SyntaxError:\n%s", e.Message)
}

type RuntimeError struct {
	Message string
	Pos     token.Position
}

func (e *RuntimeError) Error() string {
	if e.Pos.Line > 0 {
		return fmt.Sprintf("RuntimeError:\n%s\nat line %d column %d", e.Message, e.Pos.Line, e.Pos.Column)
	}
	return fmt.Sprintf("RuntimeError:\n%s", e.Message)
}

func NewSyntaxError(pos token.Position, format string, args ...interface{}) *SyntaxError {
	return &SyntaxError{
		Message: fmt.Sprintf(format, args...),
		Pos:     pos,
	}
}

func NewRuntimeError(pos token.Position, format string, args ...interface{}) *RuntimeError {
	return &RuntimeError{
		Message: fmt.Sprintf(format, args...),
		Pos:     pos,
	}
}
