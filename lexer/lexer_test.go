package lexer

import (
	"testing"

	"github.com/username/rustx/token"
)

func TestNextToken(t *testing.T) {
	input := `var five = 5
ten := 10

func add(x int, y int) int {
  return x + y
}

result := add(five, ten)
!- / * 5
5 < 10 > 5

if (5 < 10) {
	return true
} else {
	return false
}

10 == 10
10 != 9
10 <= 11
11 >= 10
5 && 5
5 || 5
x++
y--
"foobar"
"foo bar"
1.23
// single line comment
/* multi
   line
   comment */
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
		{token.DECLARE_ASSIGN, ":="},
		{token.INT, "10"},
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
		{token.IDENT, "result"},
		{token.DECLARE_ASSIGN, ":="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.RBRACE, "}"},
		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.INT, "10"},
		{token.NOT_EQ, "!="},
		{token.INT, "9"},
		{token.INT, "10"},
		{token.LTE, "<="},
		{token.INT, "11"},
		{token.INT, "11"},
		{token.GTE, ">="},
		{token.INT, "10"},
		{token.INT, "5"},
		{token.AND, "&&"},
		{token.INT, "5"},
		{token.INT, "5"},
		{token.OR, "||"},
		{token.INT, "5"},
		{token.IDENT, "x"},
		{token.INC, "++"},
		{token.IDENT, "y"},
		{token.DEC, "--"},
		{token.STRING, "foobar"},
		{token.STRING, "foo bar"},
		{token.FLOAT, "1.23"},
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
