package lexer

import (
	"rustx/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `
// Single-line comment
var five = 5;
x := 10.5;
/* Multi-line
   comment */
func add(x int, y int) int {
	return x + y;
}

if (5 < 10) {
	return true;
} else {
	return false;
}

10 == 10;
10 != 9;
10 <= 10;
10 >= 10;
x++;
x--;
!true;
a && b || c;
x % 2;
"hello\nworld";
["rustx", 42]
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

		{token.FUNCTION, "func"},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.IDENT, "int"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.IDENT, "int"},
		{token.RPAREN, ")"},
		{token.IDENT, "int"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},

		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},

		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		{token.INT, "10"},
		{token.NOT_EQ, "!="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},

		{token.INT, "10"},
		{token.LE, "<="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		{token.INT, "10"},
		{token.GE, ">="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		{token.IDENT, "x"},
		{token.PLUS_PLUS, "++"},
		{token.SEMICOLON, ";"},

		{token.IDENT, "x"},
		{token.MINUS_MINUS, "--"},
		{token.SEMICOLON, ";"},

		{token.BANG, "!"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},

		{token.IDENT, "a"},
		{token.AND, "&&"},
		{token.IDENT, "b"},
		{token.OR, "||"},
		{token.IDENT, "c"},
		{token.SEMICOLON, ";"},

		{token.IDENT, "x"},
		{token.PERCENT, "%"},
		{token.INT, "2"},
		{token.SEMICOLON, ";"},

		{token.STRING, "hello\nworld"},
		{token.SEMICOLON, ";"},

		{token.LBRACKET, "["},
		{token.STRING, "rustx"},
		{token.COMMA, ","},
		{token.INT, "42"},
		{token.RBRACKET, "]"},

		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q (literal=%q)",
				i, tt.expectedType, tok.Type, tok.Literal)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestConsecutiveComments(t *testing.T) {
	input := `
// comment 1
// comment 2
/* comment 3 */
// comment 4
x := 10
`
	l := New(input)
	tok := l.NextToken()
	if tok.Type != token.IDENT || tok.Literal != "x" {
		t.Fatalf("expected IDENT x, got %v (%s)", tok.Type, tok.Literal)
	}
}
