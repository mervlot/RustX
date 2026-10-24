package lexer

import (
	"testing"

	"rustx/token"
)

func TestNextToken(t *testing.T) {
	input := `
var x = 5;
y := 10.5;
// single line comment
/* multiline
comment */
func add(a int, b int) int {
	return a + b;
}
if x == 5 && y >= 10.0 || x != 0 {
	println("Hello \n world")
}
i++
i--
arr[0] = 1 % 2
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
		{token.DECLARE, ":="},
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
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},

		{token.IF, "if"},
		{token.IDENT, "x"},
		{token.EQ, "=="},
		{token.INT, "5"},
		{token.AND, "&&"},
		{token.IDENT, "y"},
		{token.GTE, ">="},
		{token.FLOAT, "10.0"},
		{token.OR, "||"},
		{token.IDENT, "x"},
		{token.NOT_EQ, "!="},
		{token.INT, "0"},
		{token.LBRACE, "{"},
		{token.IDENT, "println"},
		{token.LPAREN, "("},
		{token.STRING, "Hello \n world"},
		{token.RPAREN, ")"},
		{token.RBRACE, "}"},

		{token.IDENT, "i"},
		{token.INC, "++"},
		{token.IDENT, "i"},
		{token.DEC, "--"},

		{token.IDENT, "arr"},
		{token.LBRACK, "["},
		{token.INT, "0"},
		{token.RBRACK, "]"},
		{token.ASSIGN, "="},
		{token.INT, "1"},
		{token.MOD, "%"},
		{token.INT, "2"},

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
