package lexer

import (
	"testing"
	"github.com/username/rustx/token"
)

func TestNextToken(t *testing.T) {
	input := `
	// single line comment
	var x = 10;
	/* multiline
	   comment */
	name := "RustX\n";
	func add(a int, b int) int {
		return a + b;
	}
	if x >= 10 && x <= 20 || x == 0 {
		x++;
		x--;
	}
	arr := [1, 2];
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
		{token.IDENT, "name"},
		{token.COLON_ASSIGN, ":="},
		{token.STRING, "RustX\n"},
		{token.SEMICOLON, ";"},
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
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.IF, "if"},
		{token.IDENT, "x"},
		{token.GTE, ">="},
		{token.INT, "10"},
		{token.AND, "&&"},
		{token.IDENT, "x"},
		{token.LTE, "<="},
		{token.INT, "20"},
		{token.OR, "||"},
		{token.IDENT, "x"},
		{token.EQ, "=="},
		{token.INT, "0"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.INC, "++"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "x"},
		{token.DEC, "--"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.IDENT, "arr"},
		{token.COLON_ASSIGN, ":="},
		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RBRACKET, "]"},
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
