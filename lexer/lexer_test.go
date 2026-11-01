package lexer

import (
	"rustx/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `
	// Single line comment
	/* Multi
	   line
	   comment */
	var x = 10;
	name := "RustX\nStarter";
	pi := 3.14;
	func add(a int, b int) int {
		return a + b;
	}
	if (x >= 10 && x <= 20) || x != 0 {
		x++;
		x--;
	}
	while true {
		for i := 0; i < 5; i++ {
			println(i);
		}
	}
	[1, 2];
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
		{token.STRING, "RustX\nStarter"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "pi"},
		{token.COLON_ASSIGN, ":="},
		{token.FLOAT, "3.14"},
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
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.GTE, ">="},
		{token.INT, "10"},
		{token.AND, "&&"},
		{token.IDENT, "x"},
		{token.LTE, "<="},
		{token.INT, "20"},
		{token.RPAREN, ")"},
		{token.OR, "||"},
		{token.IDENT, "x"},
		{token.NOT_EQ, "!="},
		{token.INT, "0"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS_PLUS, "++"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "x"},
		{token.MINUS_MINUS, "--"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.WHILE, "while"},
		{token.TRUE, "true"},
		{token.LBRACE, "{"},
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
		{token.IDENT, "println"},
		{token.LPAREN, "("},
		{token.IDENT, "i"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.RBRACE, "}"},
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
