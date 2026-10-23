package lexer

import (
	"testing"

	"rustx/token"
)

func TestNextToken(t *testing.T) {
	input := `
var age = 16
name := "Muhammad"
func add(a int, b int) int {
    return a + b
}
if age >= 18 && name == "Muhammad" || false {
    x := 3.14 % 2
    i++
    j--
}
// single line comment
/* multi
line
comment */
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

		{token.FUNC, "func"},
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

		{token.RETURN, "return"},
		{token.IDENT, "a"},
		{token.PLUS, "+"},
		{token.IDENT, "b"},

		{token.RBRACE, "}"},

		{token.IF, "if"},
		{token.IDENT, "age"},
		{token.GTE, ">="},
		{token.INT, "18"},
		{token.AND, "&&"},
		{token.IDENT, "name"},
		{token.EQ, "=="},
		{token.STRING, "Muhammad"},
		{token.OR, "||"},
		{token.FALSE, "false"},
		{token.LBRACE, "{"},

		{token.IDENT, "x"},
		{token.COLON_ASSIGN, ":="},
		{token.FLOAT, "3.14"},
		{token.PERCENT, "%"},
		{token.INT, "2"},

		{token.IDENT, "i"},
		{token.INC, "++"},

		{token.IDENT, "j"},
		{token.DEC, "--"},

		{token.RBRACE, "}"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
