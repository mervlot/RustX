package lexer

import (
	"testing"

	"github.com/username/rustx/token"
)

func TestNextToken(t *testing.T) {
	input := `var age = 16
name := "Muhammad\n"
func add(a int, b float) int {
	// comment
	/* multiline
	   comment */
	if a >= 10 && b <= 20.5 {
		return a + b
	}
	a++
	b--
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
		{token.FUNC, "func"},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "a"},
		{token.IDENT, "int"},
		{token.COMMA, ","},
		{token.IDENT, "b"},
		{token.IDENT, "float"},
		{token.RPAREN, ")"},
		{token.IDENT, "int"},
		{token.LBRACE, "{"},
		{token.IF, "if"},
		{token.IDENT, "a"},
		{token.GTE, ">="},
		{token.INT, "10"},
		{token.AND, "&&"},
		{token.IDENT, "b"},
		{token.LTE, "<="},
		{token.FLOAT, "20.5"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.IDENT, "a"},
		{token.PLUS, "+"},
		{token.IDENT, "b"},
		{token.RBRACE, "}"},
		{token.IDENT, "a"},
		{token.INC, "++"},
		{token.IDENT, "b"},
		{token.DEC, "--"},
		{token.RETURN, "return"},
		{token.INT, "0"},
		{token.RBRACE, "}"},
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
