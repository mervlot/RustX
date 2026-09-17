package lexer

import (
	"testing"

	"github.com/username/rustx/token"
)

func TestNextTokenBasic(t *testing.T) {
	input := `var x = 5;
name := "Muhammad\n"
func add(a int, b int) float {
    return a + b
}
if (x >= 5 && x <= 10) || x != 0 {
    x++
    x--
}
/* multiline
comment */
// single line comment
y := 3.14
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
		{token.IDENT, "float"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.IDENT, "a"},
		{token.PLUS, "+"},
		{token.IDENT, "b"},
		{token.RBRACE, "}"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.GTE, ">="},
		{token.INT, "5"},
		{token.AND, "&&"},
		{token.IDENT, "x"},
		{token.LTE, "<="},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.OR, "||"},
		{token.IDENT, "x"},
		{token.NOT_EQ, "!="},
		{token.INT, "0"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.INC, "++"},
		{token.IDENT, "x"},
		{token.DEC, "--"},
		{token.RBRACE, "}"},
		{token.IDENT, "y"},
		{token.COLON_ASSIGN, ":="},
		{token.FLOAT, "3.14"},
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
