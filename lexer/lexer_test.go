package lexer

import (
	"testing"
	"github.com/username/rustx/token"
)

func TestNextToken(t *testing.T) {
	input := `var five = 5;
x := 10.5;

// Single line comment
/* Multi-line
   comment */

func add(x int, y int) int {
	return x + y;
}

if 5 < 10 && 10 >= 5 || 1 != 2 {
	println("Hello \n world!");
}

x++;
y--;
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.VAR, "var"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.IDENT, "x"},
		{token.COLON_ASSIGN, ":="},
		{token.FLOAT, "10.5"},
		{token.SEMICOLON, ";"},

		{token.FUNC, "func"},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.IDENT, "int"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.IDENT, "int"},
		{token.RPAREN, ")"},
		{token.IDENT, "int"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},

		{token.IF, "if"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.AND, "&&"},
		{token.INT, "10"},
		{token.GT_EQ, ">="},
		{token.INT, "5"},
		{token.OR, "||"},
		{token.INT, "1"},
		{token.NOT_EQ, "!="},
		{token.INT, "2"},
		{token.LBRACE, "{"},
		{token.IDENT, "println"},
		{token.LPAREN, "("},
		{token.STRING, "Hello \n world!"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},

		{token.IDENT, "x"},
		{token.INC, "++"},
		{token.SEMICOLON, ";"},

		{token.IDENT, "y"},
		{token.DEC, "--"},
		{token.SEMICOLON, ";"},

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
