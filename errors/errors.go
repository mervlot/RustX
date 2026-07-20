package errors

import (
	"fmt"
)

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
