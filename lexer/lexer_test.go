package lexer

import (
	"testing"

	"github.com/username/rustx/token"
)

func TestNextToken(t *testing.T) {
	input := `var five = 5;
x := 10.5;
// single line comment
/* multiline
   comment */
func add(x int, y int) int {
	return x + y;
}
if (x != 5 && y >= 10 || z <= 20) {
	x++;
	y--;
}
string_val := "hello\nworld";
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
		{token.COLON_EQ, ":="},
		{token.FLOAT, "10.5"},
		{token.SEMICOLON, ";"},
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
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.NOT_EQ, "!="},
		{token.INT, "5"},
		{token.AND, "&&"},
		{token.IDENT, "y"},
		{token.GTE, ">="},
		{token.INT, "10"},
		{token.OR, "||"},
		{token.IDENT, "z"},
		{token.LTE, "<="},
		{token.INT, "20"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.INC, "++"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "y"},
		{token.DEC, "--"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.IDENT, "string_val"},
		{token.COLON_EQ, ":="},
		{token.STRING, "hello\nworld"},
		{token.SEMICOLON, ";"},
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
