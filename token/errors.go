package token

import "github.com/username/rustx/errors"

func (t Token) SyntaxError(msg string) *errors.SyntaxError {
	return &errors.SyntaxError{
		Message: msg,
		Line:    t.Line,
		Column:  t.Column,
	}
}
