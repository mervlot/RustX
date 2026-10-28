package errors

import "fmt"

type SyntaxError struct {
	Message string
	Line    int
	Column  int
}

func (e *SyntaxError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf("SyntaxError:\n%s\nat line %d column %d", e.Message, e.Line, e.Column)
	}
	return fmt.Sprintf("SyntaxError: %s", e.Message)
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
	return fmt.Sprintf("RuntimeError: %s", e.Message)
}
