package errors

import (
	"fmt"
)

type SyntaxError struct {
	Message string
	Line    int
	Column  int
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("SyntaxError: %s at line %d column %d", e.Message, e.Line, e.Column)
}

type RuntimeError struct {
	Message string
	Line    int
	Column  int
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("RuntimeError: %s at line %d column %d", e.Message, e.Line, e.Column)
}

func NewSyntaxError(msg string, line, col int) *SyntaxError {
	return &SyntaxError{Message: msg, Line: line, Column: col}
}

func NewRuntimeError(msg string, line, col int) *RuntimeError {
	return &RuntimeError{Message: msg, Line: line, Column: col}
}
