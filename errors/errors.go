package errors

import "fmt"

type SyntaxError struct {
	Line    int
	Column  int
	Message string
}

func (e *SyntaxError) Error() string {
	if e.Line > 0 && e.Column > 0 {
		return fmt.Sprintf("SyntaxError: %s at line %d column %d", e.Message, e.Line, e.Column)
	}
	return fmt.Sprintf("SyntaxError: %s", e.Message)
}

type RuntimeError struct {
	Line    int
	Column  int
	Message string
}

func (e *RuntimeError) Error() string {
	if e.Line > 0 && e.Column > 0 {
		return fmt.Sprintf("RuntimeError: %s at line %d column %d", e.Message, e.Line, e.Column)
	}
	return fmt.Sprintf("RuntimeError: %s", e.Message)
}
