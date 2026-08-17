package lexer

import (
	"testing"

	"github.com/username/rustx/token"
)

func TestNextToken(t *testing.T) {
	input := `var five = 5;
x := 10.5;
// This is a single line comment
/* This is a
   multiline comment */
func add(a int, b int) int {
    return a + b;
}

if (x >= 10.0 && five <= 5) || five != 0 {
    println("Hello, \"world\"!\n");
}
for i := 0; i < 5; i++ {
    i--;
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
		{token.SEMICOLON, ";"},

		{token.IDENT, "x"},
		{token.COLON_ASSIGN, ":="},
		{token.FLOAT, "10.5"},
		{token.SEMICOLON, ";"},

		{token.FUNC, "func"},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "a"},
		{token.INT_TYPE, "int"},
		{token.COMMA, ","},
		{token.IDENT, "b"},
		{token.INT_TYPE, "int"},
		{token.RPAREN, ")"},
		{token.INT_TYPE, "int"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.IDENT, "a"},
		{token.PLUS, "+"},
		{token.IDENT, "b"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},

		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.GE, ">="},
		{token.FLOAT, "10.0"},
		{token.AND, "&&"},
		{token.IDENT, "five"},
		{token.LE, "<="},
		{token.INT, "5"},
		{token.RPAREN, ")"},
		{token.OR, "||"},
		{token.IDENT, "five"},
		{token.NOT_EQ, "!="},
		{token.INT, "0"},
		{token.LBRACE, "{"},
		{token.IDENT, "println"},
		{token.LPAREN, "("},
		{token.STRING, "Hello, \"world\"!\n"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},

		{token.FOR, "for"},
		{token.IDENT, "i"},
		{token.COLON_ASSIGN, ":="},
		{token.INT, "0"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "i"},
		{token.LT, "<"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "i"},
		{token.INC, "++"},
		{token.LBRACE, "{"},
		{token.IDENT, "i"},
		{token.DEC, "--"},
		{token.SEMICOLON, ";"},
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
