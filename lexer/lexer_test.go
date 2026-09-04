package lexer

import (
	"testing"

	"github.com/username/rustx/token"
)

func TestNextToken(t *testing.T) {
	input := `var x = 5;
y := 10.5
/* multi-line
   comment */
// single line comment
func add(a int, b int) int {
	return a + b
}
if x == 5 && y >= 10.0 || x != 0 {
	println("hello\nworld")
}
for i := 0; i < 10; i++ {
	i--
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
		{token.FUNC, "func"},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "a"},
		{token.TYPE_INT, "int"},
		{token.COMMA, ","},
		{token.IDENT, "b"},
		{token.TYPE_INT, "int"},
		{token.RPAREN, ")"},
		{token.TYPE_INT, "int"},
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
		{token.GTE, ">="},
		{token.FLOAT, "10.0"},
		{token.OR, "||"},
		{token.IDENT, "x"},
		{token.NOT_EQ, "!="},
		{token.INT, "0"},
		{token.LBRACE, "{"},
		{token.IDENT, "println"},
		{token.LPAREN, "("},
		{token.STRING, "hello\nworld"},
		{token.RPAREN, ")"},
		{token.RBRACE, "}"},
		{token.FOR, "for"},
		{token.IDENT, "i"},
		{token.COLON_ASSIGN, ":="},
		{token.INT, "0"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "i"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "i"},
		{token.INC, "++"},
		{token.LBRACE, "{"},
		{token.IDENT, "i"},
		{token.DEC, "--"},
		{token.RBRACE, "}"},
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
