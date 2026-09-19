package errors

import (
	"fmt"

	"github.com/username/rustx/token"
)

// SyntaxError represents a syntax error encountered during lexing or parsing.
type SyntaxError struct {
	Message string
	Line    int
	Column  int
}

func (e *SyntaxError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf("SyntaxError: %s at line %d column %d", e.Message, e.Line, e.Column)
	}
	return fmt.Sprintf("SyntaxError: %s", e.Message)
}

func NewSyntaxError(tok token.Token, msg string) *SyntaxError {
	return &SyntaxError{
		Message: msg,
		Line:    tok.Line,
		Column:  tok.Column,
	}
}

// RuntimeError represents an error encountered during script evaluation.
type RuntimeError struct {
	Message string
	Line    int
	Column  int
}

func (e *RuntimeError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf("RuntimeError: %s at line %d column %d", e.Message, e.Line, e.Column)
	}
	return fmt.Sprintf("RuntimeError: %s", e.Message)
}

func NewRuntimeError(tok token.Token, msg string) *RuntimeError {
	return &RuntimeError{
		Message: msg,
		Line:    tok.Line,
		Column:  tok.Column,
	}
}
