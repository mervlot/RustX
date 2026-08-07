package errors

import "fmt"

// SyntaxError represents a syntax or lexical analysis error with detailed location information.
type SyntaxError struct {
	Message string
	File    string
	Line    int
	Column  int
}

func (e *SyntaxError) Error() string {
	fileStr := ""
	if e.File != "" {
		fileStr = fmt.Sprintf(" in %s", e.File)
	}
	return fmt.Sprintf("SyntaxError: %s%s at line %d, column %d", e.Message, fileStr, e.Line, e.Column)
}

// RuntimeError represents an error occurring during program execution/evaluation.
type RuntimeError struct {
	Message string
	File    string
	Line    int
	Column  int
}

func (e *RuntimeError) Error() string {
	fileStr := ""
	if e.File != "" {
		fileStr = fmt.Sprintf(" in %s", e.File)
	}
	return fmt.Sprintf("RuntimeError: %s%s at line %d, column %d", e.Message, fileStr, e.Line, e.Column)
}
