package errors

import "fmt"

type SyntaxError struct {
	Message string
	Line    int
	Column  int
	File    string
}

func (e *SyntaxError) Error() string {
	if e.File != "" {
		return fmt.Sprintf("SyntaxError: %s at %s:%d:%d", e.Message, e.File, e.Line, e.Column)
	}
	return fmt.Sprintf("SyntaxError: %s at line %d column %d", e.Message, e.Line, e.Column)
}

type RuntimeError struct {
	Message string
	Line    int
	Column  int
	File    string
}

func (e *RuntimeError) Error() string {
	if e.File != "" {
		return fmt.Sprintf("RuntimeError: %s at %s:%d:%d", e.Message, e.File, e.Line, e.Column)
	}
	return fmt.Sprintf("RuntimeError: %s at line %d column %d", e.Message, e.Line, e.Column)
}
