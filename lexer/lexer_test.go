package lexer

import (
	"rustx/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `
// Single line comment
var x = 5;
y := 10.5;
/* Multi line
   comment */
func add(a, b) {
	return a + b;
}

if x < y && y >= 10 {
	x++
	y--
} else if x == y || x != y {
	println("hello\n\"world\"")
}

for i := 0; i < 5; i++ {
	while true {
		break
	}
}

x % 2
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
		{token.SEMICOLON, ";"},
		{token.FUNC, "func"},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "a"},
		{token.COMMA, ","},
		{token.IDENT, "b"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.IDENT, "a"},
		{token.PLUS, "+"},
		{token.IDENT, "b"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.IF, "if"},
		{token.IDENT, "x"},
		{token.LT, "<"},
		{token.IDENT, "y"},
		{token.AND, "&&"},
		{token.IDENT, "y"},
		{token.GT_EQ, ">="},
		{token.INT, "10"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS_PLUS, "++"},
		{token.IDENT, "y"},
		{token.MINUS_MINUS, "--"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.IF, "if"},
		{token.IDENT, "x"},
		{token.EQ, "=="},
		{token.IDENT, "y"},
		{token.OR, "||"},
		{token.IDENT, "x"},
		{token.NOT_EQ, "!="},
		{token.IDENT, "y"},
		{token.LBRACE, "{"},
		{token.IDENT, "println"},
		{token.LPAREN, "("},
		{token.STRING, "hello\n\"world\""},
		{token.RPAREN, ")"},
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
		{token.PLUS_PLUS, "++"},
		{token.LBRACE, "{"},
		{token.WHILE, "while"},
		{token.TRUE, "true"},
		{token.LBRACE, "{"},
		{token.IDENT, "break"},
		{token.RBRACE, "}"},
		{token.RBRACE, "}"},
		{token.IDENT, "x"},
		{token.PERCENT, "%"},
		{token.INT, "2"},
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
