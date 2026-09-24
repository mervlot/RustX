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

func NewSyntaxError(line, column int, message string, args ...interface{}) *SyntaxError {
	formattedMsg := message
	if len(args) > 0 {
		formattedMsg = fmt.Sprintf(message, args...)
	}
	return &SyntaxError{
		Message: formattedMsg,
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

func NewRuntimeError(line, column int, message string, args ...interface{}) *RuntimeError {
	formattedMsg := message
	if len(args) > 0 {
		formattedMsg = fmt.Sprintf(message, args...)
	}
	return &RuntimeError{
		Message: formattedMsg,
		Line:    line,
		Column:  column,
	}
}
