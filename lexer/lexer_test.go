package lexer

import (
	"testing"

	"rustx/token"
)

func TestNextToken(t *testing.T) {
	input := `var age = 16
name := "Muhammad"
// Single line comment
/*
  Multi line
  comment
*/
func add(a int, b int) int {
	if a == b {
		return a + b
	}
	return a - b
}
x := 3.14
count++
count--
x % 2
a && b || c
[1, 2]
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.VAR, "var"},
		{token.IDENT, "age"},
		{token.ASSIGN, "="},
		{token.INT, "16"},

		{token.IDENT, "name"},
		{token.COLON_ASSIGN, ":="},
		{token.STRING, "Muhammad"},

		{token.FUNCTION, "func"},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "a"},
		{token.IDENT, "int"},
		{token.COMMA, ","},
		{token.IDENT, "b"},
		{token.IDENT, "int"},
		{token.RPAREN, ")"},
		{token.IDENT, "int"},
		{token.LBRACE, "{"},

		{token.IF, "if"},
		{token.IDENT, "a"},
		{token.EQ, "=="},
		{token.IDENT, "b"},
		{token.LBRACE, "{"},

		{token.RETURN, "return"},
		{token.IDENT, "a"},
		{token.PLUS, "+"},
		{token.IDENT, "b"},

		{token.RBRACE, "}"},

		{token.RETURN, "return"},
		{token.IDENT, "a"},
		{token.MINUS, "-"},
		{token.IDENT, "b"},

		{token.RBRACE, "}"},

		{token.IDENT, "x"},
		{token.COLON_ASSIGN, ":="},
		{token.FLOAT, "3.14"},

		{token.IDENT, "count"},
		{token.INC, "++"},

		{token.IDENT, "count"},
		{token.DEC, "--"},

		{token.IDENT, "x"},
		{token.MODULO, "%"},
		{token.INT, "2"},

		{token.IDENT, "a"},
		{token.AND, "&&"},
		{token.IDENT, "b"},
		{token.OR, "||"},
		{token.IDENT, "c"},

		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RBRACKET, "]"},

		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q (literal %q)",
				i, tt.expectedType, tok.Type, tok.Literal)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
