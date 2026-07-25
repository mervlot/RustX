package lexer

import (
	"testing"

	"github.com/username/rustx/token"
)

func TestNextToken(t *testing.T) {
	input := `var five = 5;
let ten = 10.5;

// This is a comment
/* This is a
   multiline comment */
/* Comment 1 */ // Comment 2
"hello\nworld"["brackets"]

func add(x, y) {
	x + y;
}

x := 5;
x++;
x--;
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

		{token.IDENT, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.FLOAT, "10.5"},
		{token.SEMICOLON, ";"},

		{token.STRING, "hello\nworld"},
		{token.LBRACKET, "["},
		{token.STRING, "brackets"},
		{token.RBRACKET, "]"},

		{token.FUNC, "func"},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},

		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},

		{token.RBRACE, "}"},

		{token.IDENT, "x"},
		{token.COLON_ASSIGN, ":="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.IDENT, "x"},
		{token.INCREMENT, "++"},
		{token.SEMICOLON, ";"},

		{token.IDENT, "x"},
		{token.DECREMENT, "--"},
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
