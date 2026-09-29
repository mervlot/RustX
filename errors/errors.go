package errors

import "fmt"

type SyntaxError struct {
	Message string
	Line    int
	Column  int
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("SyntaxError: %s\n  at line %d column %d", e.Message, e.Line, e.Column)
}

func NewSyntaxError(line, column int, format string, a ...interface{}) *SyntaxError {
	return &SyntaxError{
		Message: fmt.Sprintf(format, a...),
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
	return fmt.Sprintf("RuntimeError: %s\n  at line %d column %d", e.Message, e.Line, e.Column)
}

func NewRuntimeError(line, column int, format string, a ...interface{}) *RuntimeError {
	return &RuntimeError{
		Message: fmt.Sprintf(format, a...),
		Line:    line,
		Column:  column,
	}
}
