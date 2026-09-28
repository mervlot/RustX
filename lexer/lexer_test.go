package lexer

import (
	"testing"

	"github.com/username/rustx/token"
)

func TestNextToken(t *testing.T) {
	input := `var age = 16
name := "Muhammad\n"
floatVal := 3.14
// single line comment
/* multiline
   comment */
func add(a int, b int) int {
	if a == b && a != 0 || b <= 10 {
		a++
		b--
		return a + b
	}
	return 0
}
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
		{token.STRING, "Muhammad\n"},

		{token.IDENT, "floatVal"},
		{token.COLON_ASSIGN, ":="},
		{token.FLOAT, "3.14"},

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

		{token.IF, "if"},
		{token.IDENT, "a"},
		{token.EQ, "=="},
		{token.IDENT, "b"},
		{token.AND, "&&"},
		{token.IDENT, "a"},
		{token.NOT_EQ, "!="},
		{token.INT, "0"},
		{token.OR, "||"},
		{token.IDENT, "b"},
		{token.LT_EQ, "<="},
		{token.INT, "10"},
		{token.LBRACE, "{"},

		{token.IDENT, "a"},
		{token.INC, "++"},
		{token.IDENT, "b"},
		{token.DEC, "--"},
		{token.RETURN, "return"},
		{token.IDENT, "a"},
		{token.PLUS, "+"},
		{token.IDENT, "b"},
		{token.RBRACE, "}"},

		{token.RETURN, "return"},
		{token.INT, "0"},
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
