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
	if e.Line > 0 && e.Column > 0 {
		return fmt.Sprintf("SyntaxError:\n%s\nat line %d column %d", e.Message, e.Line, e.Column)
	}
	return fmt.Sprintf("SyntaxError:\n%s", e.Message)
}

func NewSyntaxError(line, col int, format string, args ...interface{}) *SyntaxError {
	msg := fmt.Sprintf(format, args...)
	return &SyntaxError{
		Message: msg,
		Line:    line,
		Column:  col,
	}
}

type RuntimeError struct {
	Message string
	Line    int
	Column  int
}

func (e *RuntimeError) Error() string {
	if e.Line > 0 && e.Column > 0 {
		return fmt.Sprintf("RuntimeError:\n%s\nat line %d column %d", e.Message, e.Line, e.Column)
	}
	return fmt.Sprintf("RuntimeError:\n%s", e.Message)
}

func NewRuntimeError(line, col int, format string, args ...interface{}) *RuntimeError {
	msg := fmt.Sprintf(format, args...)
	return &RuntimeError{
		Message: msg,
		Line:    line,
		Column:  col,
	}
}
