package ast

import (
	"testing"
	"github.com/username/rustx/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&VarStatement{
				Token: token.Token{Type: token.VAR, Literal: "var"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	expected := "var myVar = anotherVar;"
	if program.String() != expected {
		t.Errorf("program.String() wrong. got=%q, expected=%q", program.String(), expected)
	}
}
