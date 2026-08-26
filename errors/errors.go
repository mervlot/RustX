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
	return fmt.Sprintf("SyntaxError: %s at line %d column %d", e.Message, e.Position.Line, e.Position.Column)
}

func NewSyntaxError(pos token.Position, msg string, args ...interface{}) *SyntaxError {
	formattedMsg := msg
	if len(args) > 0 {
		formattedMsg = fmt.Sprintf(msg, args...)
	}
	return &SyntaxError{
		Message:  formattedMsg,
		Position: pos,
	}
}

type RuntimeError struct {
	Message  string
	Position token.Position
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("RuntimeError: %s at line %d column %d", e.Message, e.Position.Line, e.Position.Column)
}

func NewRuntimeError(pos token.Position, msg string, args ...interface{}) *RuntimeError {
	formattedMsg := msg
	if len(args) > 0 {
		formattedMsg = fmt.Sprintf(msg, args...)
	}
	return &RuntimeError{
		Message:  formattedMsg,
		Position: pos,
	}
}
