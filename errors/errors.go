package errors

import (
	"fmt"
	"github.com/username/rustx/token"
)

type SyntaxError struct {
	Message  string
	Position token.Position
	File     string
}

func (e *SyntaxError) Error() string {
	fileStr := ""
	if e.File != "" {
		fileStr = fmt.Sprintf(" in %s", e.File)
	}
	return fmt.Sprintf("SyntaxError: %s%s at line %d column %d", e.Message, fileStr, e.Position.Line, e.Position.Column)
}

type RuntimeError struct {
	Message  string
	Position token.Position
	File     string
}

func (e *RuntimeError) Error() string {
	fileStr := ""
	if e.File != "" {
		fileStr = fmt.Sprintf(" in %s", e.File)
	}
	return fmt.Sprintf("RuntimeError: %s%s at line %d column %d", e.Message, fileStr, e.Position.Line, e.Position.Column)
}

func NewSyntaxError(pos token.Position, message string) *SyntaxError {
	return &SyntaxError{
		Message:  message,
		Position: pos,
	}
}

func NewRuntimeError(pos token.Position, message string) *RuntimeError {
	return &RuntimeError{
		Message:  message,
		Position: pos,
	}
}
