package lexer

import (
	"testing"
	"github.com/username/rustx/token"
)

func TestNextToken(t *testing.T) {
	input := `var five = 5
ten := 10.5

// Single line comment
/* Multi line
   comment */
func add(x int, y int) int {
	return x + y
}

if 10 >= 5 && 5 <= 10 {
	println("hello\nworld")
} else if x == 1 || y != 2 {
	x++
	y--
}
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.VAR, "var"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},

		{token.IDENT, "ten"},
		{token.COLON_ASSIGN, ":="},
		{token.FLOAT, "10.5"},

		{token.FUNCTION, "func"},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.INT_TYPE, "int"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.INT_TYPE, "int"},
		{token.RPAREN, ")"},
		{token.INT_TYPE, "int"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.RBRACE, "}"},

		{token.IF, "if"},
		{token.INT, "10"},
		{token.GTE, ">="},
		{token.INT, "5"},
		{token.AND, "&&"},
		{token.INT, "5"},
		{token.LTE, "<="},
		{token.INT, "10"},
		{token.LBRACE, "{"},
		{token.IDENT, "println"},
		{token.LPAREN, "("},
		{token.STRING, "hello\nworld"},
		{token.RPAREN, ")"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.IF, "if"},
		{token.IDENT, "x"},
		{token.EQ, "=="},
		{token.INT, "1"},
		{token.OR, "||"},
		{token.IDENT, "y"},
		{token.NOT_EQ, "!="},
		{token.INT, "2"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.INC, "++"},
		{token.IDENT, "y"},
		{token.DEC, "--"},
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
