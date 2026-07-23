package lexer

import (
	"testing"

	"github.com/username/rustx/token"
)

func TestNextToken(t *testing.T) {
	input := `var x = 5;
y := 12.34;
/* This is a multi-line
   comment */
// and a single-line comment
func add(a, b) {
	return a + b;
}
"hello\nworld"[0]
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedLine    int
		expectedColumn  int
	}{
		{token.VAR, "var", 1, 1},
		{token.IDENT, "x", 1, 5},
		{token.ASSIGN, "=", 1, 7},
		{token.INT, "5", 1, 9},
		{token.SEMICOLON, ";", 1, 10},

		{token.IDENT, "y", 2, 1},
		{token.COLON_ASSIGN, ":=", 2, 3},
		{token.FLOAT, "12.34", 2, 6},
		{token.SEMICOLON, ";", 2, 11},

		// comments on lines 3, 4, 5 are skipped
		{token.FUNC, "func", 6, 1},
		{token.IDENT, "add", 6, 6},
		{token.LPAREN, "(", 6, 9},
		{token.IDENT, "a", 6, 10},
		{token.COMMA, ",", 6, 11},
		{token.IDENT, "b", 6, 13},
		{token.RPAREN, ")", 6, 14},
		{token.LBRACE, "{", 6, 16},

		{token.RETURN, "return", 7, 2},
		{token.IDENT, "a", 7, 9},
		{token.PLUS, "+", 7, 11},
		{token.IDENT, "b", 7, 13},
		{token.SEMICOLON, ";", 7, 14},

		{token.RBRACE, "}", 8, 1},

		{token.STRING, "hello\nworld", 9, 1},
		{token.LBRACKET, "[", 9, 15},
		{token.INT, "0", 9, 16},
		{token.RBRACKET, "]", 9, 17},

		{token.EOF, "", 10, 1},
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

		if tok.Line != tt.expectedLine {
			t.Fatalf("tests[%d] - line wrong. expected=%d, got=%d",
				i, tt.expectedLine, tok.Line)
		}

		if tok.Column != tt.expectedColumn {
			t.Fatalf("tests[%d] - column wrong. expected=%d, got=%d",
				i, tt.expectedColumn, tok.Column)
		}
	}
}
