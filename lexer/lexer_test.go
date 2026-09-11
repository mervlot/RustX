package lexer

import (
	"testing"

	"github.com/username/rustx/token"
)

func TestNextToken(t *testing.T) {
	input := `var five = 5;
x := 10.5
// single line comment
/* multi
   line
   comment */
func add(a int, b int) int {
	return a + b
}
if (x >= 10.5 && five != 0) || false {
	println("hello\nworld")
}
i++
j--
arr[0]
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
		{token.GTE, ">="},
		{token.FLOAT, "10.5"},
		{token.AND, "&&"},
		{token.IDENT, "five"},
		{token.NOT_EQ, "!="},
		{token.INT, "0"},
		{token.RPAREN, ")"},
		{token.OR, "||"},
		{token.FALSE, "false"},
		{token.LBRACE, "{"},
		{token.IDENT, "println"},
		{token.LPAREN, "("},
		{token.STRING, "hello\nworld"},
		{token.RPAREN, ")"},
		{token.RBRACE, "}"},
		{token.IDENT, "i"},
		{token.INC, "++"},
		{token.IDENT, "j"},
		{token.DEC, "--"},
		{token.IDENT, "arr"},
		{token.LBRACK, "["},
		{token.INT, "0"},
		{token.RBRACK, "]"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q (literal=%q, line=%d, col=%d)",
				i, tt.expectedType, tok.Type, tok.Literal, tok.Line, tok.Column)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
