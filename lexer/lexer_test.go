package lexer

import (
	"testing"

	"github.com/username/rustx/token"
)

func TestNextToken(t *testing.T) {
	input := `var five = 5;
x := 10;
y := 3.14;
// single line comment
/* multiline
   comment */
func add(x int, y int) int {
	return x + y;
}

if (5 < 10) {
	return true;
} else {
	return false;
}

10 == 10;
10 != 9;
"hello world"
"escape\nline"
i++
i--
x && y || z
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
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		{token.IDENT, "y"},
		{token.COLON_EQ, ":="},
		{token.FLOAT, "3.14"},
		{token.SEMICOLON, ";"},

		{token.FUNC, "func"},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.TYPE_INT, "int"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.TYPE_INT, "int"},
		{token.RPAREN, ")"},
		{token.TYPE_INT, "int"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},

		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},

		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		{token.INT, "10"},
		{token.NOT_EQ, "!="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},

		{token.STRING, "hello world"},
		{token.STRING, "escape\nline"},

		{token.IDENT, "i"},
		{token.INC, "++"},
		{token.IDENT, "i"},
		{token.DEC, "--"},

		{token.IDENT, "x"},
		{token.AND, "&&"},
		{token.IDENT, "y"},
		{token.OR, "||"},
		{token.IDENT, "z"},

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
