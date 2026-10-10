package lexer

import (
	"testing"

	"github.com/username/rustx/token"
)

func TestNextToken(t *testing.T) {
	input := `var five = 5;
x := 10.5;

// Single line comment
/*
  Multi line comment
*/

func add(a int, b int) int {
	return a + b;
}

if (five < 10) {
	println("Hello\nWorld");
} else {
	print("Bye");
}

five++;
x--;

true && false || !true
a == b
a != b
a <= b
a >= b
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
		{token.SEMICOLON, ";"},

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
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},

		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "println"},
		{token.LPAREN, "("},
		{token.STRING, "Hello\nWorld"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.IDENT, "print"},
		{token.LPAREN, "("},
		{token.STRING, "Bye"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},

		{token.IDENT, "five"},
		{token.INCREMENT, "++"},
		{token.SEMICOLON, ";"},

		{token.IDENT, "x"},
		{token.DECREMENT, "--"},
		{token.SEMICOLON, ";"},

		{token.TRUE, "true"},
		{token.AND, "&&"},
		{token.FALSE, "false"},
		{token.OR, "||"},
		{token.BANG, "!"},
		{token.TRUE, "true"},

		{token.IDENT, "a"},
		{token.EQ, "=="},
		{token.IDENT, "b"},

		{token.IDENT, "a"},
		{token.NOT_EQ, "!="},
		{token.IDENT, "b"},

		{token.IDENT, "a"},
		{token.LTE, "<="},
		{token.IDENT, "b"},

		{token.IDENT, "a"},
		{token.GTE, ">="},
		{token.IDENT, "b"},

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
