package lexer

import (
	"testing"

	"github.com/username/rustx/token"
)

func TestNextToken(t *testing.T) {
	input := `var five = 5
name := "RustX"

/* multi-line
   comment */
func add(x int, y int) int {
	return x + y // single-line comment
}

if 5 < 10 {
	println("yes")
} else {
	println("no")
}

val := 3.14
a := true
b := false
i++
i--
a && b || true
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.VAR, "var"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},

		{token.IDENT, "name"},
		{token.COLON_ASSIGN, ":="},
		{token.STRING, "RustX"},

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
		{token.RBRACE, "}"},

		{token.IF, "if"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.LBRACE, "{"},
		{token.IDENT, "println"},
		{token.LPAREN, "("},
		{token.STRING, "yes"},
		{token.RPAREN, ")"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.IDENT, "println"},
		{token.LPAREN, "("},
		{token.STRING, "no"},
		{token.RPAREN, ")"},
		{token.RBRACE, "}"},

		{token.IDENT, "val"},
		{token.COLON_ASSIGN, ":="},
		{token.FLOAT, "3.14"},

		{token.IDENT, "a"},
		{token.COLON_ASSIGN, ":="},
		{token.TRUE, "true"},

		{token.IDENT, "b"},
		{token.COLON_ASSIGN, ":="},
		{token.FALSE, "false"},

		{token.IDENT, "i"},
		{token.INC, "++"},

		{token.IDENT, "i"},
		{token.DEC, "--"},

		{token.IDENT, "a"},
		{token.AND, "&&"},
		{token.IDENT, "b"},
		{token.OR, "||"},
		{token.TRUE, "true"},

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
