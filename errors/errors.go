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
		return fmt.Sprintf("SyntaxError: %s at line %d column %d", e.Message, e.Position.Line, e.Position.Column)
	}
	return fmt.Sprintf("SyntaxError: %s", e.Message)
}

type RuntimeError struct {
	Message  string
	Position token.Position
}

func (e *RuntimeError) Error() string {
	if e.Position.Line > 0 {
		return fmt.Sprintf("RuntimeError: %s at line %d column %d", e.Message, e.Position.Line, e.Position.Column)
	}
	return fmt.Sprintf("RuntimeError: %s", e.Message)
}
