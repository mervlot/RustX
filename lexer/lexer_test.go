package lexer

import (
	"testing"
	"github.com/username/rustx/token"
)

func TestNextToken(t *testing.T) {
	input := `var five = 5
name := "RustX"
func add(x int, y int) int {
	return x + y
}
// comment
/* multiline
   comment */
if (five == 5) {
	five++
	five--
} else if 10 != 5 {
	five = five * 2 / 1
}
var f = 3.14
a && b || c
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.VAR, "var"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},

		{token.IDENT, "name"},
		{token.COLON_ASSIGN, ":="},
		{token.STRING, "RustX"},

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
		{token.RBRACE, "}"},

		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.EQ, "=="},
		{token.INT, "5"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "five"},
		{token.INC, "++"},
		{token.IDENT, "five"},
		{token.DEC, "--"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.IF, "if"},
		{token.INT, "10"},
		{token.NOT_EQ, "!="},
		{token.INT, "5"},
		{token.LBRACE, "{"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.IDENT, "five"},
		{token.ASTERISK, "*"},
		{token.INT, "2"},
		{token.SLASH, "/"},
		{token.INT, "1"},
		{token.RBRACE, "}"},

		{token.VAR, "var"},
		{token.IDENT, "f"},
		{token.ASSIGN, "="},
		{token.FLOAT, "3.14"},

		{token.IDENT, "a"},
		{token.AND, "&&"},
		{token.IDENT, "b"},
		{token.OR, "||"},
		{token.IDENT, "c"},

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
