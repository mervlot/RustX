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

func NewSyntaxError(tok token.Token, format string, a ...interface{}) *SyntaxError {
	msg := format
	if len(a) > 0 {
		msg = fmt.Sprintf(format, a...)
	}
	return &SyntaxError{
		Message: msg,
		Token:   tok,
	}
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

func NewRuntimeError(tok token.Token, format string, a ...interface{}) *RuntimeError {
	msg := format
	if len(a) > 0 {
		msg = fmt.Sprintf(format, a...)
	}
	return &RuntimeError{
		Message: msg,
		Token:   tok,
	}
}
