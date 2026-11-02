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

func NewSyntaxError(msg string, line int, col int) *SyntaxError {
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
	if e.Line > 0 {
		return fmt.Sprintf("RuntimeError:\n%s\nat line %d column %d", e.Message, e.Line, e.Column)
	}
	return fmt.Sprintf("RuntimeError:\n%s", e.Message)
}

func NewRuntimeError(msg string, line int, col int) *RuntimeError {
	return &RuntimeError{
		Message: msg,
		Line:    line,
		Column:  col,
	}
}
