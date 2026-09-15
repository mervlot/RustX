package errors

import "fmt"

// SyntaxError represents a compilation or parsing error with line and column information.
type SyntaxError struct {
	Message string
	Line    int
	Column  int
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("SyntaxError: %s at line %d column %d", e.Message, e.Line, e.Column)
}

// RuntimeError represents an error during program evaluation/execution.
type RuntimeError struct {
	Message string
	Line    int
	Column  int
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("RuntimeError: %s at line %d column %d", e.Message, e.Line, e.Column)
}

// NewSyntaxError creates a new SyntaxError.
func NewSyntaxError(line, col int, format string, args ...interface{}) *SyntaxError {
	return &SyntaxError{
		Message: fmt.Sprintf(format, args...),
		Line:    line,
		Column:  col,
	}
}

// NewRuntimeError creates a new RuntimeError.
func NewRuntimeError(line, col int, format string, args ...interface{}) *RuntimeError {
	return &RuntimeError{
		Message: fmt.Sprintf(format, args...),
		Line:    line,
		Column:  col,
	}
}
