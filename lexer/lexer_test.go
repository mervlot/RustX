package lexer

import (
	"testing"

	"github.com/username/rustx/token"
)

func TestNextToken(t *testing.T) {
	input := `var five = 5;
x := 10.5;
func add(x int, y int) int {
	return x + y;
}

// single line comment
/* multiline
   comment */
if (x == 10.5 && five != 0 || x <= 20) {
	x++;
	five--;
}
string_val := "hello\nworld \"quoted\"";
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
		{token.EQ, "=="},
		{token.FLOAT, "10.5"},
		{token.AND, "&&"},
		{token.IDENT, "five"},
		{token.NOT_EQ, "!="},
		{token.INT, "0"},
		{token.OR, "||"},
		{token.IDENT, "x"},
		{token.LTE, "<="},
		{token.INT, "20"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.INC, "++"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "five"},
		{token.DEC, "--"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.IDENT, "string_val"},
		{token.COLON_EQ, ":="},
		{token.STRING, "hello\nworld \"quoted\""},
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
