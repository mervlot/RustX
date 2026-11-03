package lexer

import (
	"rustx/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `
// Sample RustX code
var x = 10;
y := 20.5;

/* Multiline
   comment */
func add(a int, b int) int {
	return a + b;
}

if x == 10 && y != 0.0 || !false {
	x++;
	y--;
} else {
	x = x % 3;
}

while (x > 0) {
	println("hello\t\"world\"\n");
}

[1, 2, "three"];
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.VAR, "var"},
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "y"},
		{token.COLON_ASSIGN, ":="},
		{token.FLOAT, "20.5"},
		{token.SEMICOLON, ";"},
		{token.FUNCTION, "func"},
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
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.IF, "if"},
		{token.IDENT, "x"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.AND, "&&"},
		{token.IDENT, "y"},
		{token.NOT_EQ, "!="},
		{token.FLOAT, "0.0"},
		{token.OR, "||"},
		{token.BANG, "!"},
		{token.FALSE, "false"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS_PLUS, "++"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "y"},
		{token.MINUS_MINUS, "--"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.IDENT, "x"},
		{token.PERCENT, "%"},
		{token.INT, "3"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.WHILE, "while"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.GT, ">"},
		{token.INT, "0"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "println"},
		{token.LPAREN, "("},
		{token.STRING, "hello\t\"world\"\n"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.COMMA, ","},
		{token.STRING, "three"},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("test[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("test[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
