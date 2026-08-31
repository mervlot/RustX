package lexer

import (
	"testing"

	"github.com/username/rustx/token"
)

func TestNextToken(t *testing.T) {
	input := `var age = 16
// single line comment
/* multiline
   comment */
name := "Muhammad\n"
func add(a int, b int) int {
    return a + b
}
if age == 16 && name != "" {
    age++
} else {
    age--
}
3.14 <= 4.0 || true != false
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
		{token.RETURN, "return"},
		{token.IDENT, "a"},
		{token.PLUS, "+"},
		{token.IDENT, "b"},
		{token.RBRACE, "}"},
		{token.IF, "if"},
		{token.IDENT, "age"},
		{token.EQ, "=="},
		{token.INT, "16"},
		{token.AND, "&&"},
		{token.IDENT, "name"},
		{token.NOT_EQ, "!="},
		{token.STRING, ""},
		{token.LBRACE, "{"},
		{token.IDENT, "age"},
		{token.INC, "++"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.IDENT, "age"},
		{token.DEC, "--"},
		{token.RBRACE, "}"},
		{token.FLOAT, "3.14"},
		{token.LTE, "<="},
		{token.FLOAT, "4.0"},
		{token.OR, "||"},
		{token.TRUE, "true"},
		{token.NOT_EQ, "!="},
		{token.FALSE, "false"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q (literal=%q)",
				i, tt.expectedType, tok.Type, tok.Literal)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
