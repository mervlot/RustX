package lexer

import (
	"testing"

	"github.com/username/rustx/token"
)

func TestNextToken(t *testing.T) {
	input := `
	var x = 5
	y := 10.5
	func add(a, b) {
		return a + b
	}
	// single line comment
	/* multi
	   line comment */
	if x == 5 && y != 0.0 || x <= 10 {
		x++
		y--
	}
	str := "hello\nworld"
	`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.VAR, "var"},
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.INT, "5"},

		{token.IDENT, "y"},
		{token.COLON_ASSIGN, ":="},
		{token.FLOAT, "10.5"},

		{token.FUNCTION, "func"},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "a"},
		{token.COMMA, ","},
		{token.IDENT, "b"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.IDENT, "a"},
		{token.PLUS, "+"},
		{token.IDENT, "b"},
		{token.RBRACE, "}"},

		{token.IF, "if"},
		{token.IDENT, "x"},
		{token.EQ, "=="},
		{token.INT, "5"},
		{token.AND, "&&"},
		{token.IDENT, "y"},
		{token.NOT_EQ, "!="},
		{token.FLOAT, "0.0"},
		{token.OR, "||"},
		{token.IDENT, "x"},
		{token.LTE, "<="},
		{token.INT, "10"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS_PLUS, "++"},
		{token.IDENT, "y"},
		{token.MINUS_MINUS, "--"},
		{token.RBRACE, "}"},

		{token.IDENT, "str"},
		{token.COLON_ASSIGN, ":="},
		{token.STRING, "hello\nworld"},

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
