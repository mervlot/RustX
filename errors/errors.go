package errors

import "fmt"

type SyntaxError struct {
	Message string
	Line    int
	Column  int
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("SyntaxError: %s at line %d column %d", e.Message, e.Line, e.Column)
}

func NewSyntaxError(line, column int, format string, args ...interface{}) *SyntaxError {
	msg := fmt.Sprintf(format, args...)
	return &SyntaxError{
		Message: msg,
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
	return fmt.Sprintf("RuntimeError: %s at line %d column %d", e.Message, e.Line, e.Column)
}

func NewRuntimeError(line, column int, format string, args ...interface{}) *RuntimeError {
	msg := fmt.Sprintf(format, args...)
	return &RuntimeError{
		Message: msg,
		Line:    line,
		Column:  column,
	}
}
