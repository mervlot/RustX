package lexer

import (
	"testing"

	"github.com/username/rustx/token"
)

func TestNextToken(t *testing.T) {
	input := `
	// This is a comment
	var x = 5;
	y := 10.5;
	/* Multi
	   line
	   comment */
	func add(a int, b int) int {
		return a + b
	}

	if (x == 5 && y != 0.0) || name == "hello\nworld" {
		x++
		y--
	}
	`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.VAR, "var"},
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "y"},
		{token.COLON_ASSIGN, ":="},
		{token.FLOAT, "10.5"},
		{token.SEMICOLON, ";"},
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
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.EQ, "=="},
		{token.INT, "5"},
		{token.AND, "&&"},
		{token.IDENT, "y"},
		{token.NOT_EQ, "!="},
		{token.FLOAT, "0.0"},
		{token.RPAREN, ")"},
		{token.OR, "||"},
		{token.IDENT, "name"},
		{token.EQ, "=="},
		{token.STRING, "hello\nworld"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.INC, "++"},
		{token.IDENT, "y"},
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
