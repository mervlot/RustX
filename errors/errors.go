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
	return fmt.Sprintf("SyntaxError: %s at line %d column %d", e.Message, e.Line, e.Col)
}

type RuntimeError struct {
	Message string
	Line    int
	Col     int
}

func (e *RuntimeError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf("RuntimeError: %s at line %d column %d", e.Message, e.Line, e.Col)
	}
	return fmt.Sprintf("RuntimeError: %s", e.Message)
}
