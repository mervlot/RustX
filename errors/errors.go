package errors

import (
	"fmt"
	"github.com/username/rustx/token"
)

type SyntaxError struct {
	Message string
	Token   token.Token
}

func (e *SyntaxError) Error() string {
	if e.Token.Line > 0 {
		return fmt.Sprintf("SyntaxError: %s at line %d column %d", e.Message, e.Token.Line, e.Token.Column)
	}
	return fmt.Sprintf("SyntaxError: %s", e.Message)
}

type RuntimeError struct {
	Message string
	Token   token.Token
}

func (e *RuntimeError) Error() string {
	if e.Token.Line > 0 {
		return fmt.Sprintf("RuntimeError: %s at line %d column %d", e.Message, e.Token.Line, e.Token.Column)
	}
	return fmt.Sprintf("RuntimeError: %s", e.Message)
}

func NewSyntaxError(tok token.Token, msg string) *SyntaxError {
	return &SyntaxError{
		Message: msg,
		Token:   tok,
	}
}

func NewRuntimeError(tok token.Token, msg string) *RuntimeError {
	return &RuntimeError{
		Message: msg,
		Token:   tok,
	}
}
