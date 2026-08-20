package lexer

import (
	"testing"

	"github.com/username/rustx/token"
)

func TestNextToken(t *testing.T) {
	input := `var age = 16
name := "Muhammad\n"
// single line comment
/* multi
   line */
func add(a int, b int) int {
    if a == b { return a + b }
    return a * b / 2.5
}
i++
i--
a && b || c <= 10 >= 5 != 0
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
		{token.ASTERISK, "*"},
		{token.IDENT, "b"},
		{token.SLASH, "/"},
		{token.FLOAT, "2.5"},
		{token.RBRACE, "}"},
		{token.IDENT, "i"},
		{token.INC, "++"},
		{token.IDENT, "i"},
		{token.DEC, "--"},
		{token.IDENT, "a"},
		{token.AND, "&&"},
		{token.IDENT, "b"},
		{token.OR, "||"},
		{token.IDENT, "c"},
		{token.LTE, "<="},
		{token.INT, "10"},
		{token.GTE, ">="},
		{token.INT, "5"},
		{token.NOT_EQ, "!="},
		{token.INT, "0"},
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
