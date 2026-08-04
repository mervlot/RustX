package lexer

import (
	"testing"

	"github.com/username/rustx/token"
)

func TestNextToken(t *testing.T) {
	input := `var x = 5;
y := 12.34;
// single comment
/* multi-line
   comment */
func add(a int, b int) int {
    return a + b;
}

"hello \t world\n"
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedLine    int
		expectedCol     int
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

		{token.FUNC, "func", 6, 1},
		{token.IDENT, "add", 6, 6},
		{token.LPAREN, "(", 6, 9},
		{token.IDENT, "a", 6, 10},
		{token.IDENT, "int", 6, 12},
		{token.COMMA, ",", 6, 15},
		{token.IDENT, "b", 6, 17},
		{token.IDENT, "int", 6, 19},
		{token.RPAREN, ")", 6, 22},
		{token.IDENT, "int", 6, 24},
		{token.LBRACE, "{", 6, 28},

		{token.RETURN, "return", 7, 5},
		{token.IDENT, "a", 7, 12},
		{token.PLUS, "+", 7, 14},
		{token.IDENT, "b", 7, 16},
		{token.SEMICOLON, ";", 7, 17},

		{token.RBRACE, "}", 8, 1},

		{token.STRING, "hello \t world\n", 10, 1},
		{token.EOF, "", 11, 1},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}

		if tok.Line != tt.expectedLine {
			t.Fatalf("tests[%d] - line wrong. expected=%d, got=%d", i, tt.expectedLine, tok.Line)
		}

		if tok.Col != tt.expectedCol {
			t.Fatalf("tests[%d] - column wrong. expected=%d, got=%d", i, tt.expectedCol, tok.Col)
		}
	}
}

func TestAdjacentBracketAndComments(t *testing.T) {
	input := `"hello"[0] // comment
// another comment
/* multiline */[1]
`
	l := New(input)

	expected := []struct {
		typ token.TokenType
		lit string
	}{
		{token.STRING, "hello"},
		{token.LBRACKET, "["},
		{token.INT, "0"},
		{token.RBRACKET, "]"},
		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.RBRACKET, "]"},
		{token.EOF, ""},
	}

	for i, tt := range expected {
		tok := l.NextToken()
		if tok.Type != tt.typ {
			t.Fatalf("test[%d] - typ wrong, expected=%q, got=%q", i, tt.typ, tok.Type)
		}
		if tok.Literal != tt.lit {
			t.Fatalf("test[%d] - lit wrong, expected=%q, got=%q", i, tt.lit, tok.Literal)
		}
	}
}
