package errors

import "fmt"

type SyntaxError struct {
	Message string
	Line    int
	Column  int
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("SyntaxError:\n%s\nat line %d column %d", e.Message, e.Line, e.Column)
}

func NewSyntaxError(message string, line, column int) *SyntaxError {
	return &SyntaxError{
		Message: message,
		Line:    line,
		Column:  column,
	}
}

type RuntimeError struct {
	Message string
	Line    int
	Column  int
}

func (e *RuntimeError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf("RuntimeError:\n%s\nat line %d column %d", e.Message, e.Line, e.Column)
	}
	return fmt.Sprintf("RuntimeError:\n%s", e.Message)
}

func NewRuntimeError(message string, line, column int) *RuntimeError {
	return &RuntimeError{
		Message: message,
		Line:    line,
		Column:  column,
	}
}
