package lexer

import (
	"testing"

	"github.com/username/rustx/token"
)

func TestNextToken(t *testing.T) {
	input := `
var x = 5;
y := 10.5;
func add(a int, b int) int {
	// single line comment
	/* multiline
	   comment */
	if (a == b) {
		return a + b;
	} else if (a != b && a <= b || a >= b) {
		a++;
		b--;
	}
	return a;
}
str := "hello \"world\"\n";
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

		{token.FUNCTION, "func"},
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

		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.IDENT, "a"},
		{token.EQ, "=="},
		{token.IDENT, "b"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},

		{token.RETURN, "return"},
		{token.IDENT, "a"},
		{token.PLUS, "+"},
		{token.IDENT, "b"},
		{token.SEMICOLON, ";"},

		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.IDENT, "a"},
		{token.NOT_EQ, "!="},
		{token.IDENT, "b"},
		{token.AND, "&&"},
		{token.IDENT, "a"},
		{token.LTE, "<="},
		{token.IDENT, "b"},
		{token.OR, "||"},
		{token.IDENT, "a"},
		{token.GTE, ">="},
		{token.IDENT, "b"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},

		{token.IDENT, "a"},
		{token.INCREMENT, "++"},
		{token.SEMICOLON, ";"},

		{token.IDENT, "b"},
		{token.DECREMENT, "--"},
		{token.SEMICOLON, ";"},

		{token.RBRACE, "}"},

		{token.RETURN, "return"},
		{token.IDENT, "a"},
		{token.SEMICOLON, ";"},

		{token.RBRACE, "}"},

		{token.IDENT, "str"},
		{token.COLON_ASSIGN, ":="},
		{token.STRING, "hello \"world\"\n"},
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
