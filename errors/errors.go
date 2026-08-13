package errors

import "fmt"

type SyntaxError struct {
	Message string
	Line    int
	Col     int
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("SyntaxError:\n%s\nat line %d column %d", e.Message, e.Line, e.Col)
}

type RuntimeError struct {
	Message string
	Line    int
	Col     int
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("RuntimeError:\n%s\nat line %d column %d", e.Message, e.Line, e.Col)
}

func NewSyntaxError(msg string, line, col int) *SyntaxError {
	return &SyntaxError{Message: msg, Line: line, Col: col}
}

func NewRuntimeError(msg string, line, col int) *RuntimeError {
	return &RuntimeError{Message: msg, Line: line, Col: col}
}
