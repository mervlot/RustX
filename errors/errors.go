package errors

import "fmt"

// SyntaxError represents an error during lexing or parsing.
type SyntaxError struct {
	Message string
	Line    int
	Column  int
	File    string
}

func (e *SyntaxError) Error() string {
	fileInfo := ""
	if e.File != "" {
		fileInfo = fmt.Sprintf(" in %s", e.File)
	}
	return fmt.Sprintf("SyntaxError: %s%s at line %d column %d", e.Message, fileInfo, e.Line, e.Column)
}

// RuntimeError represents an error during evaluation.
type RuntimeError struct {
	Message string
	Line    int
	Column  int
	File    string
}

func (e *RuntimeError) Error() string {
	fileInfo := ""
	if e.File != "" {
		fileInfo = fmt.Sprintf(" in %s", e.File)
	}
	return fmt.Sprintf("RuntimeError: %s%s at line %d column %d", e.Message, fileInfo, e.Line, e.Column)
}
