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

type RuntimeError struct {
	Message string
	Line    int
	Column  int
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("RuntimeError: %s at line %d column %d", e.Message, e.Line, e.Column)
}

func NewSyntaxError(line, col int, format string, args ...interface{}) *SyntaxError {
	return &SyntaxError{
		Message: fmt.Sprintf(format, args...),
		Line:    line,
		Column:  col,
	}
}

func NewRuntimeError(line, col int, format string, args ...interface{}) *RuntimeError {
	return &RuntimeError{
		Message: fmt.Sprintf(format, args...),
		Line:    line,
		Column:  col,
	}
}
