package errors

import "fmt"

type SyntaxError struct {
	Message string
	Line    int
	Column  int
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("SyntaxError: %s\nat line %d column %d", e.Message, e.Line, e.Column)
}

func NewSyntaxError(line, col int, message string) *SyntaxError {
	return &SyntaxError{
		Message: message,
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
	return fmt.Sprintf("RuntimeError: %s\nat line %d column %d", e.Message, e.Line, e.Column)
}

func NewRuntimeError(line, col int, message string) *RuntimeError {
	return &RuntimeError{
		Message: message,
		Line:    line,
		Column:  col,
	}
}
